package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"gorm.io/gorm"
)

func (s *storage) CreateUser(username, passwordHash string, familyID uuid.UUID) (*models.User, error) {
	user := models.User{}
	user.ID = uuid.New()
	user.Username = username
	user.PasswordHash = passwordHash
	user.FamilyID = familyID
	user.StartDate = time.Now()
	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	return &user, nil
}

func (s *storage) GetUser(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(StorageError, err)
	}
	return &user, nil
}

func (s *storage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(StorageError, err)
	}
	return &user, nil
}

func (s *storage) PutUser(user *models.User) error {
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}
	return nil
}

func (s *storage) GetFamilyByName(name string) (*models.Family, error) {
	var family models.Family
	if err := s.db.Where("name = ?", name).First(&family).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(StorageError, err)
	}
	return &family, nil
}

func (s *storage) CreateFamily(name string) (*models.Family, error) {
	family := models.Family{}
	family.ID = uuid.New()
	family.Name = name
	if err := s.db.Create(&family).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	return &family, nil
}

func (s *storage) GetAllFamilyIDs() ([]uuid.UUID, error) {
	var families []models.Family
	if err := s.db.Find(&families).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	ids := make([]uuid.UUID, len(families))
	for i, f := range families {
		ids[i] = f.ID
	}
	return ids, nil
}
