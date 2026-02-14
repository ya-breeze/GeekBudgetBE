package database

import (
	"errors"
	"fmt"

	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"gorm.io/gorm"
)

func (s *storage) CreateImage(data []byte, contentType string) (models.Image, error) {
	image := models.Image{
		Data:        data,
		ContentType: contentType,
	}

	if err := s.db.Create(&image).Error; err != nil {
		return models.Image{}, fmt.Errorf(StorageError, err)
	}

	return image, nil
}

func (s *storage) GetImage(id string) (models.Image, error) {
	var image models.Image
	if err := s.db.Where("id = ?", id).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Image{}, ErrNotFound
		}
		return models.Image{}, fmt.Errorf(StorageError, err)
	}

	return image, nil
}

func (s *storage) DeleteImage(id string) error {
	if err := s.db.Where("id = ?", id).Delete(&models.Image{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}
	return nil
}
