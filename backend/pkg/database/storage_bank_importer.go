package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) GetBankImporters(familyID uuid.UUID) ([]goserver.BankImporter, error) {
	result, err := s.db.Model(&models.BankImporter{}).Where("family_id = ?", familyID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	importers := make([]goserver.BankImporter, 0)
	for result.Next() {
		var imp models.BankImporter
		if err := s.db.ScanRows(result, &imp); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		importers = append(importers, imp.FromDB())
	}

	return importers, nil
}

func (s *storage) CreateBankImporter(familyID uuid.UUID, bankImporter *goserver.BankImporterNoId,
) (goserver.BankImporter, error) {
	data := models.BankImporterToDB(bankImporter, familyID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, familyID, "BankImporter", data.ID.String(), "CREATED", nil, data); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	s.log.Info("BankImporter created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) UpdateBankImporter(familyID uuid.UUID, id string, bankImporter goserver.BankImporterNoIdInterface,
) (goserver.BankImporter, error) {
	return performUpdate[models.BankImporter, goserver.BankImporterNoIdInterface, goserver.BankImporter](s, familyID, "BankImporter", id, bankImporter,
		models.BankImporterToDB,
		func(m *models.BankImporter) goserver.BankImporter { return m.FromDB() },
		func(m *models.BankImporter, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteBankImporter(familyID uuid.UUID, id string) error {
	var data models.BankImporter
	if err := s.db.Where("id = ? AND family_id = ?", id, familyID).First(&data).Error; err == nil {
		if err := s.recordAuditLog(s.db, familyID, "BankImporter", id, "DELETED", &data, nil); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}
	}

	if err := s.db.Where("id = ? AND family_id = ?", id, familyID).Delete(&models.BankImporter{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetBankImporter(familyID uuid.UUID, id string) (goserver.BankImporter, error) {
	var data models.BankImporter
	if err := s.db.Where("id = ? AND family_id = ?", id, familyID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BankImporter{}, ErrNotFound
		}

		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) GetAllBankImporters() ([]ImportInfo, error) {
	result, err := s.db.Model(&models.BankImporter{}).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	importers := make([]ImportInfo, 0)
	for result.Next() {
		var imp models.BankImporter
		if err := s.db.ScanRows(result, &imp); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		importers = append(importers, ImportInfo{
			FamilyID:         imp.FamilyID,
			BankImporterID:   imp.ID.String(),
			BankImporterType: imp.Type,
			FetchAll:         imp.FetchAll,
		})
	}

	return importers, nil
}

func (s *storage) GetBankImporterFiles(familyID uuid.UUID) ([]goserver.BankImporterFile, error) {
	var files []models.BankImporterFile
	if err := s.db.Where("family_id = ?", familyID).Order("upload_date DESC").Find(&files).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.BankImporterFile, len(files))
	for i, f := range files {
		result[i] = f.FromDB()
	}

	return result, nil
}

func (s *storage) GetBankImporterFile(familyID uuid.UUID, id string) (models.BankImporterFile, error) {
	var file models.BankImporterFile
	if err := s.db.Where("family_id = ? AND id = ?", familyID, id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BankImporterFile{}, ErrNotFound
		}
		return models.BankImporterFile{}, fmt.Errorf(StorageError, err)
	}

	return file, nil
}

func (s *storage) CreateBankImporterFile(familyID uuid.UUID, file *models.BankImporterFile) (goserver.BankImporterFile, error) {
	file.FamilyID = familyID
	file.ID = uuid.New()
	if file.UploadDate.IsZero() {
		file.UploadDate = time.Now()
	}

	if err := s.db.Create(file).Error; err != nil {
		return goserver.BankImporterFile{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, familyID, "BankImporterFile", file.ID.String(), "CREATED", nil, file); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return file.FromDB(), nil
}

func (s *storage) DeleteBankImporterFile(familyID uuid.UUID, id string) error {
	var data models.BankImporterFile
	if err := s.db.Where("id = ? AND family_id = ?", id, familyID).First(&data).Error; err == nil {
		if err := s.recordAuditLog(s.db, familyID, "BankImporterFile", id, "DELETED", &data, nil); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}
	}

	if err := s.db.Where("family_id = ? AND id = ?", familyID, id).Delete(&models.BankImporterFile{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}
	return nil
}
