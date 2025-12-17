package dto

import "github.com/google/uuid"

type BalanceResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt string    `json:"last_updated_at"`
}