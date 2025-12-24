package dto

import "github.com/google/uuid"

type BalanceResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt string    `json:"last_updated_at"`
}

type BalanceHistoryItem struct {
	Action         string   `json:"action"`
	PreviousAmount float64  `json:"previous_amount"`
	NewAmount      float64  `json:"new_amount"`
	ChangeAmount   float64  `json:"change_amount"`
	RelatedUserID  *string  `json:"related_user_id,omitempty"`
	TransactionID  *string  `json:"transaction_id,omitempty"`
	CreatedAt      string   `json:"created_at"`
}
type BalanceAtTimeRequest struct {
	Timestamp int64 `json:"timestamp" validate:"required,gt=0"`
}

type BalanceAtTimeResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	AsOf      string    `json:"as_of"`    
	IsExact   bool      `json:"is_exact"`
}