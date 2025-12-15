package repository

import (
	"backend-path/app/models"

	"gorm.io/gorm"
)

type TestRepository struct {}

func (r *TestRepository) GetAll() (tests []*models.Test, err error) {
	if err := DB.Find(&tests).Order("updated_at DESC").Error; err != nil {
		return nil, err
	}

	return tests, nil
}

func (r *TestRepository) InsertMany(tx *gorm.DB, tests []models.Test, batchSize int) error {
	if err := tx.CreateInBatches(&tests, batchSize).Error; err != nil {
		return err
	}

	return nil
}