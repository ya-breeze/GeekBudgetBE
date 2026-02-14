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

func (s *storage) GetBankImporters(userID string) ([]goserver.BankImporter, error) {
	result, err := s.db.Model(&models.BankImporter{}).Where("user_id = ?", userID).Rows()
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

func (s *storage) CreateBankImporter(userID string, bankImporter *goserver.BankImporterNoId,
) (goserver.BankImporter, error) {
	data := models.BankImporterToDB(bankImporter, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}
	s.log.Info("BankImporter created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) UpdateBankImporter(userID string, id string, bankImporter goserver.BankImporterNoIdInterface,
) (goserver.BankImporter, error) {
	return performUpdate[models.BankImporter, goserver.BankImporterNoIdInterface, goserver.BankImporter](s, userID, id, bankImporter,
		models.BankImporterToDB,
		func(m *models.BankImporter) goserver.BankImporter { return m.FromDB() },
		func(m *models.BankImporter, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteBankImporter(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.BankImporter{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetBankImporter(userID string, id string) (goserver.BankImporter, error) {
	var data models.BankImporter
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
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
			UserID:           imp.UserID,
			BankImporterID:   imp.ID.String(),
			BankImporterType: imp.Type,
			FetchAll:         imp.FetchAll,
		})
	}

	return importers, nil
}

func (s *storage) GetBankImporterFiles(userID string) ([]goserver.BankImporterFile, error) {
	var files []models.BankImporterFile
	if err := s.db.Where("user_id = ?", userID).Order("upload_date DESC").Find(&files).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.BankImporterFile, len(files))
	for i, f := range files {
		result[i] = f.FromDB()
	}

	return result, nil
}

func (s *storage) GetBankImporterFile(userID string, id string) (models.BankImporterFile, error) {
	var file models.BankImporterFile
	if err := s.db.Where("user_id = ? AND id = ?", userID, id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BankImporterFile{}, ErrNotFound
		}
		return models.BankImporterFile{}, fmt.Errorf(StorageError, err)
	}

	return file, nil
}

func (s *storage) CreateBankImporterFile(userID string, file *models.BankImporterFile) (goserver.BankImporterFile, error) {
	file.UserID = userID
	file.ID = uuid.New()
	if file.UploadDate.IsZero() {
		file.UploadDate = time.Now()
	}

	if err := s.db.Create(file).Error; err != nil {
		return goserver.BankImporterFile{}, fmt.Errorf(StorageError, err)
	}

	return file.FromDB(), nil
}

func (s *storage) DeleteBankImporterFile(userID string, id string) error {
	if err := s.db.Where("user_id = ? AND id = ?", userID, id).Delete(&models.BankImporterFile{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}
	return nil
}
