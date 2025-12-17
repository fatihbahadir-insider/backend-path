package dto

import "github.com/google/uuid"

type DepositRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,max=1000000"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,max=1000000"`
}

type TransferRequest struct {
	ToUserID string  `json:"to_user_id" validate:"required,uuid"`
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

type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
}