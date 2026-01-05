package dto

import "github.com/google/uuid"

type CreditRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,max=1000000"`
}

type DebitRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,max=1000000"`
}

type TransferRequest struct {
	ToUserID uuid.UUID  `json:"to_user_id" validate:"required,uuid"`
	Amount   float64 `json:"amount" validate:"required,gt=0,max=1000000"`
}

type TransactionResponse struct {
	ID         uuid.UUID `json:"id"`
	FromUserID *string   `json:"from_user_id,omitempty"`
	ToUserID   *string   `json:"to_user_id,omitempty"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
}

type TransactionStatsResponse struct {
	TotalProcessed   int64   `json:"total_processed"`
	TotalSuccessful  int64   `json:"total_successful"`
	TotalFailed      int64   `json:"total_failed"`
	PendingInQueue   int     `json:"pending_in_queue"`
	TotalCredited    float64 `json:"total_credited"`
	TotalDebited     float64 `json:"total_debited"`
	TotalTransferred float64 `json:"total_transferred"`
}