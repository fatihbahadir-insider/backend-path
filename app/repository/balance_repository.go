package repository

import (
	"backend-path/app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBalanceRepository interface {
	FindByUserID(userID uuid.UUID) (*models.Balance, error)
	FindByUserIDForUpdate(tx *gorm.DB, userID uuid.UUID) (*models.Balance, error)
	Create(balance *models.Balance) error
	Update(tx *gorm.DB, balance *models.Balance) error
	Upsert(tx *gorm.DB, balance *models.Balance) error
	GetBalanceHistory(userID uuid.UUID) ([]models.AuditLog, error)
	GetBalanceAtTime(userID uuid.UUID, timestamp time.Time) (*models.AuditLog, error)
}

type BalanceRepository struct {}

func NewBalanceRepository() *BalanceRepository {
	return &BalanceRepository{}
}

func (r *BalanceRepository) FindByUserID(userID uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	if err := DB.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *BalanceRepository) FindByUserIDForUpdate(tx *gorm.DB, userID uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&balance).Error
	return &balance, err
}

func (r *BalanceRepository) Create(balance *models.Balance) error {
	return DB.Create(balance).Error
}

func (r *BalanceRepository) Update(tx *gorm.DB, balance *models.Balance) error {
	return tx.Save(balance).Error
}

func (r *BalanceRepository) Upsert(tx *gorm.DB, balance *models.Balance) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount", "last_updated_at"}),
	}).Create(balance).Error
}

func (r *BalanceRepository) GetBalanceHistory(userID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := DB.Model(&models.AuditLog{}).Where("entity_type = ? AND entity_id = ?", models.EntityBalance, userID).Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *BalanceRepository) GetBalanceAtTime(userID uuid.UUID, timestamp time.Time) (*models.AuditLog, error) {
	var log models.AuditLog
	err := DB.Where("entity_type = ? AND entity_id = ? AND created_at <= ?",
		models.EntityBalance, userID, timestamp).
		Order("created_at DESC").
		First(&log).Error
	return &log, err
}