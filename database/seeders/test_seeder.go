package seeders

import (
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/utils"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestSeeder struct{}

func (s *TestSeeder) Seed(db *gorm.DB) error {
	utils.Logger.Info("✅ seed data from TestSeeder")

	maxSize := 200
	batchSize := 50

	var testsData []models.Test
	for i := 1; i <= maxSize; i++ {
		data := models.Test{
			GUID:        uuid.MustParse(gofakeit.UUID()),
			Title:       gofakeit.Sentence(5),
		}

		testsData = append(testsData, data)
	}

	testsRepo := new(repository.TestRepository)
	return testsRepo.InsertMany(db, testsData, batchSize)
}