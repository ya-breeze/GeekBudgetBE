package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"gorm.io/gorm"
)

func (s *storage) GetTransactions(userID string, dateFrom, dateTo time.Time, onlySuspicious bool) ([]goserver.Transaction, error) {
	req := s.db.Model(&models.Transaction{}).Where("user_id = ? AND merged_into_id IS NULL", userID)
	if onlySuspicious {
		// Filter transactions where suspicious_reasons is not null and not an empty JSON array
		req = req.Where("suspicious_reasons IS NOT NULL AND suspicious_reasons != '[]' AND suspicious_reasons != ''")
	}
	if !dateFrom.IsZero() {
		req = req.Where("date >= ?", dateFrom)
	}
	if !dateTo.IsZero() {
		req = req.Where("date < ?", dateTo)
	}
	req = req.Order("date")

	result, err := req.Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	transactions := make([]goserver.Transaction, 0)
	for result.Next() {
		var tr models.Transaction
		if err := s.db.ScanRows(result, &tr); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		transactions = append(transactions, tr.FromDB())
	}

	// Populate DuplicateTransactionIds in batch
	if len(transactions) > 0 {
		ids := make([]uuid.UUID, len(transactions))
		for i, t := range transactions {
			id, _ := uuid.Parse(t.Id)
			ids[i] = id
		}

		var relationships []models.TransactionDuplicate
		if err := s.db.Where("user_id = ? AND transaction_id1 IN ?", userID, ids).Find(&relationships).Error; err == nil {
			relMap := make(map[string][]string)
			for _, r := range relationships {
				t1 := r.TransactionID1.String()
				t2 := r.TransactionID2.String()
				relMap[t1] = append(relMap[t1], t2)
			}

			for i := range transactions {
				if dups, ok := relMap[transactions[i].Id]; ok {
					transactions[i].DuplicateTransactionIds = dups
				}
			}
		}

		// Populate MergedTransactionIds in batch (find transactions that were merged into these from archive)
		var mergedRecords []struct {
			KeptTransactionID     uuid.UUID
			OriginalTransactionID uuid.UUID
		}
		if err := s.db.Model(&models.MergedTransaction{}).
			Select("kept_transaction_id, original_transaction_id").
			Where("user_id = ? AND kept_transaction_id IN ?", userID, ids).
			Find(&mergedRecords).Error; err == nil {

			mergedMap := make(map[string][]string)
			for _, r := range mergedRecords {
				key := r.KeptTransactionID.String()
				mergedMap[key] = append(mergedMap[key], r.OriginalTransactionID.String())
			}
			for i := range transactions {
				if merged, ok := mergedMap[transactions[i].Id]; ok {
					transactions[i].MergedTransactionIds = merged
				}
			}
		}
	}

	return transactions, nil
}

func (s *storage) GetTransactionsIncludingDeleted(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error) {
	req := s.db.Model(&models.Transaction{}).Unscoped().Where("user_id = ?", userID)
	if !dateFrom.IsZero() {
		req = req.Where("date >= ?", dateFrom)
	}
	if !dateTo.IsZero() {
		req = req.Where("date < ?", dateTo)
	}
	req = req.Order("date")

	result, err := req.Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	transactions := make([]goserver.Transaction, 0)
	for result.Next() {
		var tr models.Transaction
		if err := s.db.ScanRows(result, &tr); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		transactions = append(transactions, tr.FromDB())
	}

	return transactions, nil
}

func (s *storage) CreateTransaction(userID string, input goserver.TransactionNoIdInterface,
) (goserver.Transaction, error) {
	if err := s.validateTransaction(userID, input); err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	t := models.TransactionToDB(input, userID)
	t.ID = uuid.New()
	if err := s.db.Create(t).Error; err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Transaction", t.ID.String(), "CREATED", t); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	s.log.Info("Transaction created", "id", t.ID)

	// Invalidate reconciliation if we inserted a transaction in the past
	// We pass empty oldMovements because it's a new transaction
	s.invalidateReconciliationIfAmountsChanged(userID, []goserver.Movement{}, models.MovementsToAPI(t.Movements), t.Date)

	return t.FromDB(), nil
}

func (s *storage) UpdateTransaction(userID string, id string, input goserver.TransactionNoIdInterface,
) (goserver.Transaction, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var t *models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Transaction{}, ErrNotFound
		}

		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Transaction", t.ID.String(), "UPDATED", t); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	// Get old movements for smart invalidation
	oldMovements := models.MovementsToAPI(t.Movements)

	if err := s.validateTransaction(userID, input); err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	t = models.TransactionToDB(input, userID)
	t.ID = idUUID
	if err := s.db.Save(&t).Error; err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	// If dismissed, clear relationships
	if t.DuplicateDismissed {
		if err := s.ClearDuplicateRelationships(userID, id); err != nil {
			s.log.Error("Failed to clear duplicate relationships on dismissal", "error", err, "id", id)
		}
	} else {
		// Revalidate duplicate links in case date/amount changed
		if err := s.RevalidateDuplicateRelationships(userID, id); err != nil {
			s.log.Error("Failed to revalidate duplicate relationships", "error", err, "id", id)
		}
	}

	// Smart invalidation: only if amounts or currencies changed
	s.invalidateReconciliationIfAmountsChanged(userID, oldMovements, models.MovementsToAPI(t.Movements), t.Date)

	// Invalidate reconciliation for all affected accounts/currencies
	for _, m := range oldMovements {
		if err := s.InvalidateReconciliation(userID, m.AccountId, m.CurrencyId, t.Date); err != nil {
			s.log.Error("Failed to invalidate reconciliation for old movement", "error", err, "transaction_id", t.ID)
		}
	}
	for _, m := range models.MovementsToAPI(t.Movements) {
		if err := s.InvalidateReconciliation(userID, m.AccountId, m.CurrencyId, t.Date); err != nil {
			s.log.Error("Failed to invalidate reconciliation for new movement", "error", err, "transaction_id", t.ID)
		}
	}

	return t.FromDB(), nil
}

func (s *storage) DeleteTransaction(userID string, id string) error {
	var t models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf(StorageError, err)
	}

	if len(t.ExternalIDs) > 0 || t.UnprocessedSources != "" {
		return ErrImportedTransactionCannotBeDeleted
	}

	if err := s.recordAuditLog(s.db, userID, "Transaction", t.ID.String(), "DELETED", &t); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	if err := s.db.Delete(&t).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	// Clear duplicate relationships if any
	if err := s.ClearDuplicateRelationships(userID, id); err != nil {
		s.log.Error("Failed to clear duplicate relationships on deletion", "error", err, "id", id)
	}

	// Invalidate reconciliation for deleted movements
	s.invalidateReconciliationIfAmountsChanged(userID, models.MovementsToAPI(t.Movements), []goserver.Movement{}, t.Date)

	return nil
}

func (s *storage) MergeTransactions(userID, keepID, mergeID string) (goserver.Transaction, error) {
	kID, err := uuid.Parse(keepID)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf("invalid keep ID: %w", err)
	}
	mID, err := uuid.Parse(mergeID)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf("invalid merge ID: %w", err)
	}

	var keepT models.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var mergeT models.Transaction
		if err := tx.Where("user_id = ? AND id = ?", userID, kID).First(&keepT).Error; err != nil {
			return fmt.Errorf("failed to find keep transaction: %w", err)
		}
		if err := tx.Where("user_id = ? AND id = ?", userID, mID).First(&mergeT).Error; err != nil {
			return fmt.Errorf("failed to find merge transaction: %w", err)
		}

		// 2. Transfer external IDs
		existingIDs := make(map[string]bool)
		for _, id := range keepT.ExternalIDs {
			existingIDs[id] = true
		}
		for _, id := range mergeT.ExternalIDs {
			if !existingIDs[id] {
				keepT.ExternalIDs = append(keepT.ExternalIDs, id)
			}
		}

		// 3. Update keepT: remove suspicious reason, inherit external IDs
		newReasons := make([]string, 0)
		for _, r := range keepT.SuspiciousReasons {
			if r != "Potential duplicate from different importer" {
				newReasons = append(newReasons, r)
			}
		}
		keepT.SuspiciousReasons = newReasons
		if err := tx.Save(&keepT).Error; err != nil {
			return fmt.Errorf("failed to update keep transaction: %w", err)
		}

		// 4. Archive and hard-delete mergeT (archive is now source of truth)
		now := time.Now()

		if err := s.archiveMergedTransaction(tx, userID, &mergeT, kID, now); err != nil {
			return err
		}

		// Hard-delete the merged transaction
		if err := tx.Unscoped().Delete(&mergeT).Error; err != nil {
			return fmt.Errorf("failed to hard-delete merge transaction: %w", err)
		}

		// 5. Clear duplicate relationships for both (they are resolved now)
		// We use tx so it's atomic within our transaction
		if err := s.clearDuplicateRelationshipsWithTx(tx, userID, mID.String()); err != nil {
			return fmt.Errorf("failed to clear relationships for merged transaction: %w", err)
		}

		// Clear specific relationship between keepT and mergeT (the one where keepT was the primary)
		if err := s.clearDuplicateRelationshipsWithTx(tx, userID, kID.String()); err != nil {
			return fmt.Errorf("failed to clear relationships for keep transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	return s.GetTransaction(userID, keepID)
}

func (s *storage) DeleteDuplicateTransaction(userID string, id, duplicateID string) error {
	s.log.Info("Deleting duplicate transaction", "id", id, "duplicate_id", duplicateID)
	return s.db.Transaction(func(tx *gorm.DB) error {
		var t, duplicate models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
			s.log.Warn("Failed to find transaction", "id", id, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		if err := tx.Where("id = ? AND user_id = ?", duplicateID, userID).First(&duplicate).Error; err != nil {
			s.log.Warn("Failed to find duplicate transaction", "id", duplicateID, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		duplicate.ExternalIDs = append(duplicate.ExternalIDs, t.ExternalIDs...)
		if err := tx.Save(&duplicate).Error; err != nil {
			s.log.Warn("Failed to update duplicate transaction", "id", duplicateID, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		// Archive and hard-delete the transaction
		now := time.Now()
		duplicateIDUUID, _ := uuid.Parse(duplicateID)
		if err := s.archiveMergedTransaction(tx, userID, &t, duplicateIDUUID, now); err != nil {
			return err
		}

		// Hard-delete the transaction (archive is source of truth)
		if err := tx.Unscoped().Delete(&t).Error; err != nil {
			s.log.Warn("Failed to hard-delete merged transaction", "id", id, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		if err := s.recordAuditLog(tx, userID, "Transaction", t.ID.String(), "MERGED", &t); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}

		return nil
	})
}

func (s *storage) GetTransaction(userID string, id string) (goserver.Transaction, error) {
	var transaction models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Transaction{}, ErrNotFound
		}

		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	apiTransaction := transaction.FromDB()
	duplicateIds, err := s.GetDuplicateTransactionIDs(userID, id)
	if err == nil {
		apiTransaction.DuplicateTransactionIds = duplicateIds
	}

	// Populate MergedTransactionIds for single transaction from archive
	var mergedIds []string
	if err := s.db.Model(&models.MergedTransaction{}).
		Where("user_id = ? AND kept_transaction_id = ?", userID, id).
		Pluck("original_transaction_id", &mergedIds).Error; err == nil {
		apiTransaction.MergedTransactionIds = mergedIds
	}

	return apiTransaction, nil
}

func (s *storage) GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error) {
	var mergedModels []models.MergedTransaction
	if err := s.db.Where("user_id = ?", userID).Order("merged_at DESC").Find(&mergedModels).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.MergedTransaction, 0, len(mergedModels))
	for _, m := range mergedModels {
		var kept models.Transaction
		if err := s.db.Where("id = ? AND user_id = ?", m.KeptTransactionID, userID).First(&kept).Error; err != nil {
			s.log.Warn("Failed to find kept transaction for merged transaction", "merged_id", m.OriginalTransactionID, "kept_id", m.KeptTransactionID)
			continue
		}

		// Create a temporary transaction structure to use FromDB
		tr := models.Transaction{
			ID:                 m.OriginalTransactionID,
			UserID:             m.UserID,
			Date:               m.Date,
			Description:        m.Description,
			Place:              m.Place,
			Tags:               m.Tags,
			PartnerName:        m.PartnerName,
			PartnerAccount:     m.PartnerAccount,
			PartnerInternalID:  m.PartnerInternalID,
			Extra:              m.Extra,
			UnprocessedSources: m.UnprocessedSources,
			ExternalIDs:        m.ExternalIDs,
			Movements:          m.Movements,
			MatcherID:          m.MatcherID,
			IsAuto:             m.IsAuto,
			SuspiciousReasons:  m.SuspiciousReasons,
		}

		result = append(result, goserver.MergedTransaction{
			Transaction: tr.FromDB(),
			MergedInto:  kept.FromDB(),
			MergedAt:    m.MergedAt,
		})
	}

	return result, nil
}

func (s *storage) GetMergedTransaction(userID, originalTransactionID string) (goserver.MergedTransaction, error) {
	// 1. Find the archived transaction
	var archived models.MergedTransaction
	if err := s.db.Where("user_id = ? AND original_transaction_id = ?", userID, originalTransactionID).First(&archived).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.MergedTransaction{}, ErrNotFound
		}
		return goserver.MergedTransaction{}, fmt.Errorf(StorageError, err)
	}

	// 2. Find the kept transaction (mergedInto)
	var kept models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", archived.KeptTransactionID, userID).First(&kept).Error; err != nil {
		s.log.Warn("Failed to find kept transaction for merged transaction", "merged_id", archived.OriginalTransactionID, "kept_id", archived.KeptTransactionID)
		// We still return the merged transaction, just without the full kept transaction details if not found
	}

	// 3. Construct the response
	// Create a temporary transaction structure to use FromDB for the archived transaction
	tr := models.Transaction{
		ID:                 archived.OriginalTransactionID,
		UserID:             archived.UserID,
		Date:               archived.Date,
		Description:        archived.Description,
		Place:              archived.Place,
		Tags:               archived.Tags,
		PartnerName:        archived.PartnerName,
		PartnerAccount:     archived.PartnerAccount,
		PartnerInternalID:  archived.PartnerInternalID,
		Extra:              archived.Extra,
		UnprocessedSources: archived.UnprocessedSources,
		ExternalIDs:        archived.ExternalIDs,
		Movements:          archived.Movements,
		MatcherID:          archived.MatcherID,
		IsAuto:             archived.IsAuto,
		SuspiciousReasons:  archived.SuspiciousReasons,
	}

	return goserver.MergedTransaction{
		Transaction: tr.FromDB(),
		MergedInto:  kept.FromDB(),
		MergedAt:    archived.MergedAt,
	}, nil
}

func (s *storage) UnmergeTransaction(userID, id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Find the archived transaction
		var archived models.MergedTransaction
		if err := tx.Where("user_id = ? AND original_transaction_id = ?", userID, id).First(&archived).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("transaction %s is not merged or archive not found", id)
			}
			return fmt.Errorf(StorageError, err)
		}

		keptID := archived.KeptTransactionID.String()

		// 2. Remove external IDs from kept transaction
		var kept models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", keptID, userID).First(&kept).Error; err == nil {
			newExternalIDs := make([]string, 0)
			for _, extID := range kept.ExternalIDs {
				found := false
				for _, mergedExtID := range archived.ExternalIDs {
					if extID == mergedExtID {
						found = true
						break
					}
				}
				if !found {
					newExternalIDs = append(newExternalIDs, extID)
				}
			}
			kept.ExternalIDs = newExternalIDs
			if err := tx.Save(&kept).Error; err != nil {
				return fmt.Errorf("failed to update kept transaction: %w", err)
			}
		}

		// 3. Recreate the transaction from archive data (don't rely on soft-deleted record)
		restoredTransaction := models.Transaction{
			ID:                 archived.OriginalTransactionID,
			UserID:             archived.UserID,
			Date:               archived.Date,
			Description:        archived.Description,
			Place:              archived.Place,
			Tags:               archived.Tags,
			PartnerName:        archived.PartnerName,
			PartnerAccount:     archived.PartnerAccount,
			PartnerInternalID:  archived.PartnerInternalID,
			Extra:              archived.Extra,
			UnprocessedSources: archived.UnprocessedSources,
			ExternalIDs:        archived.ExternalIDs,
			Movements:          archived.Movements,
			MatcherID:          archived.MatcherID,
			IsAuto:             archived.IsAuto,
			SuspiciousReasons:  archived.SuspiciousReasons,
			// MergedIntoID and MergedAt are left nil (transaction is no longer merged)
		}

		// Create the restored transaction
		if err := tx.Create(&restoredTransaction).Error; err != nil {
			return fmt.Errorf("failed to recreate transaction: %w", err)
		}

		// 4. Delete from archive
		if err := tx.Delete(&archived).Error; err != nil {
			return fmt.Errorf("failed to delete from archive during unmerge: %w", err)
		}

		// 5. Record history
		if err := s.recordAuditLog(tx, userID, "Transaction", restoredTransaction.ID.String(), "UNMERGED", &restoredTransaction); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}

		return nil
	})
}

func (s *storage) GetDuplicateTransactionIDs(userID, transactionID string) ([]string, error) {
	var duplicates []models.TransactionDuplicate
	err := s.db.Where("user_id = ? AND transaction_id1 = ?", userID, transactionID).Find(&duplicates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get duplicate relationships: %w", err)
	}

	ids := make([]string, len(duplicates))
	for i, d := range duplicates {
		ids[i] = d.TransactionID2.String()
	}
	return ids, nil
}

func (s *storage) AddDuplicateRelationship(userID, transactionID1, transactionID2 string) error {
	id1, err := uuid.Parse(transactionID1)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 1: %w", err)
	}
	id2, err := uuid.Parse(transactionID2)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 2: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Add T1 -> T2
		var d1 models.TransactionDuplicate
		err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).First(&d1).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				d1 = models.TransactionDuplicate{
					UserID:         userID,
					TransactionID1: id1,
					TransactionID2: id2,
				}
				if err := tx.Create(&d1).Error; err != nil {
					return fmt.Errorf("failed to create link T1->T2: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check link T1->T2: %w", err)
			}
		}

		// Add T2 -> T1
		var d2 models.TransactionDuplicate
		err = tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).First(&d2).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				d2 = models.TransactionDuplicate{
					UserID:         userID,
					TransactionID1: id2,
					TransactionID2: id1,
				}
				if err := tx.Create(&d2).Error; err != nil {
					return fmt.Errorf("failed to create link T2->T1: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check link T2->T1: %w", err)
			}
		}

		return nil
	})
}

func (s *storage) RemoveDuplicateRelationship(userID, transactionID1, transactionID2 string) error {
	id1, err := uuid.Parse(transactionID1)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 1: %w", err)
	}
	id2, err := uuid.Parse(transactionID2)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 2: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete link T1->T2: %w", err)
		}
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete link T2->T1: %w", err)
		}
		return nil
	})
}

func (s *storage) ClearDuplicateRelationships(userID, transactionID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.clearDuplicateRelationshipsWithTx(tx, userID, transactionID)
	})
}

func (s *storage) clearDuplicateRelationshipsWithTx(tx *gorm.DB, userID, transactionID string) error {
	id, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	// 1. Find all duplicates linked to this transaction to update them later
	var duplicates []models.TransactionDuplicate
	if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Find(&duplicates).Error; err != nil {
		return fmt.Errorf("failed to find duplicate links: %w", err)
	}

	// 2. Delete bidirectional links
	for _, d := range duplicates {
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, d.TransactionID2, id).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete inverse link: %w", err)
		}
	}

	if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Delete(&models.TransactionDuplicate{}).Error; err != nil {
		return fmt.Errorf("failed to delete primary links: %w", err)
	}

	// 3. Sync suspicious reasons for all affected transactions
	// This ensures that if they no longer have duplicates, the flag is removed.
	affectedIDs := []uuid.UUID{id}
	for _, d := range duplicates {
		affectedIDs = append(affectedIDs, d.TransactionID2)
	}

	for _, affectedID := range affectedIDs {
		if err := s.syncDuplicateSuspiciousReason(tx, userID, affectedID); err != nil {
			s.log.Error("Failed to sync suspicious reason", "error", err, "id", affectedID)
		}
	}

	return nil
}

func (s *storage) syncDuplicateSuspiciousReason(tx *gorm.DB, userID string, transactionID uuid.UUID) error {
	// Check if any duplicate links remain for this transaction
	var count int64
	if err := tx.Model(&models.TransactionDuplicate{}).Where("user_id = ? AND transaction_id1 = ?", userID, transactionID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Still has duplicates, keep the reason
	}

	// No more duplicates, remove the reason if present
	var t models.Transaction
	if err := tx.Where("user_id = ? AND id = ?", userID, transactionID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Transaction might have been deleted (e.g. in Merge)
		}
		return err
	}

	newReasons := make([]string, 0)
	reasonsChanged := false
	for _, r := range t.SuspiciousReasons {
		if r == models.DuplicateReason {
			reasonsChanged = true
			continue
		}
		newReasons = append(newReasons, r)
	}

	if reasonsChanged {
		if err := tx.Model(&t).Update("suspicious_reasons", newReasons).Error; err != nil {
			return fmt.Errorf("failed to update suspicious reasons: %w", err)
		}
	}

	return nil
}

// RevalidateDuplicateRelationships re-checks all duplicate links for a transaction.
// If a linked transaction no longer passes IsDuplicate, the link is removed and
// suspicious reasons are synchronized for both transactions.
func (s *storage) RevalidateDuplicateRelationships(userID, transactionID string) error {
	id, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get the transaction
		var t models.Transaction
		if err := tx.Where("user_id = ? AND id = ?", userID, id).First(&t).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil // Transaction doesn't exist, nothing to revalidate
			}
			return err
		}

		// 2. Get all linked duplicates
		var duplicates []models.TransactionDuplicate
		if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Find(&duplicates).Error; err != nil {
			return err
		}

		// 3. For each link, re-check IsDuplicate
		for _, d := range duplicates {
			var linkedT models.Transaction
			if err := tx.Where("user_id = ? AND id = ?", userID, d.TransactionID2).First(&linkedT).Error; err != nil {
				// Linked transaction doesn't exist, clean up the link
				s.removeDuplicateLinkWithTx(tx, userID, id, d.TransactionID2)
				continue
			}

			if !utils.IsDuplicate(t.Date, t.Movements, linkedT.Date, linkedT.Movements) {
				// No longer duplicates, remove bidirectional link
				s.removeDuplicateLinkWithTx(tx, userID, id, d.TransactionID2)
				// Sync suspicious reasons for both
				s.syncDuplicateSuspiciousReason(tx, userID, id)
				s.syncDuplicateSuspiciousReason(tx, userID, d.TransactionID2)
			}
		}

		return nil
	})
}

func (s *storage) removeDuplicateLinkWithTx(tx *gorm.DB, userID string, id1, id2 uuid.UUID) {
	tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).Delete(&models.TransactionDuplicate{})
	tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).Delete(&models.TransactionDuplicate{})
}

func (s *storage) archiveMergedTransaction(tx *gorm.DB, userID string,
	merged *models.Transaction, keptID uuid.UUID, mergedAt time.Time,
) error {
	archive := models.MergedTransaction{
		ID:                    uuid.New(),
		UserID:                userID,
		KeptTransactionID:     keptID,
		OriginalTransactionID: merged.ID,
		Date:                  merged.Date,
		Description:           merged.Description,
		Place:                 merged.Place,
		Tags:                  merged.Tags,
		PartnerName:           merged.PartnerName,
		PartnerAccount:        merged.PartnerAccount,
		PartnerInternalID:     merged.PartnerInternalID,
		Extra:                 merged.Extra,
		UnprocessedSources:    merged.UnprocessedSources,
		ExternalIDs:           merged.ExternalIDs,
		Movements:             merged.Movements,
		MatcherID:             merged.MatcherID,
		IsAuto:                merged.IsAuto,
		SuspiciousReasons:     merged.SuspiciousReasons,
		MergedAt:              mergedAt,
	}

	if err := tx.Create(&archive).Error; err != nil {
		return fmt.Errorf("failed to create merged transaction archive: %w", err)
	}

	return nil
}

func (s *storage) validateTransaction(userID string, transaction goserver.TransactionNoIdInterface) error {
	for _, m := range transaction.GetMovements() {
		// 1. Validate Amount (decimal.Decimal doesn't have NaN/Inf, but we can check if it's uninitialized if needed)
		// No specific check for decimal.Decimal yet, as it's inherently more stable than float64

		// 2. Validate AccountId
		if m.AccountId != "" {
			var count int64
			if err := s.db.Model(&models.Account{}).Where("user_id = ? AND id = ?", userID, m.AccountId).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("account %s not found", m.AccountId)
			}
		}

		// 3. Validate CurrencyId
		if m.CurrencyId != "" {
			var count int64
			if err := s.db.Model(&models.Currency{}).Where("user_id = ? AND id = ?", userID, m.CurrencyId).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("currency %s not found", m.CurrencyId)
			}
		} else {
			return errors.New("currency ID is required")
		}
	}
	return nil
}
