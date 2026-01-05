package repository

import (
	"backend-path/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(tx *gorm.DB, transaction *models.Transaction) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error)
	Update(tx *gorm.DB, transaction *models.Transaction) error
	GetDB() *gorm.DB
}

type TransactionRepository struct {}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) Create(tx *gorm.DB, transaction *models.Transaction) error {
	if tx == nil {
		tx = DB
	}

	return tx.Create(transaction).Error
}

func (r *TransactionRepository) FindByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := DB.Preload("FromUser").Preload("ToUser").Where("id = ?", id).First(&transaction).Error

	return &transaction, err
}

func (r *TransactionRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := DB.Model(&models.Transaction{}).Where("from_user_id = ? OR to_user_id = ?", userID, userID)

	query.Count(&total)

	err := query.Preload("FromUser").Preload("ToUser").
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&transactions).Error

	return transactions, total, err
}

func (r *TransactionRepository) Update(tx *gorm.DB, transaction *models.Transaction) error {
	if tx == nil {
		tx = DB
	}

	return tx.Save(transaction).Error
}

func (r *TransactionRepository) GetDB() *gorm.DB {
	return DB
}