package models

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	UserID        uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Amount        float64   `json:"amount" gorm:"type:decimal(15,2);default:0"`
	LastUpdatedAt time.Time `json:"last_updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (Balance) TableName() string {
	return "balances"
}