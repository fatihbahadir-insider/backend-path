package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType uint

const (
	TxTypeDeposit  TransactionType = iota + 1 
	TxTypeWithdraw                            
	TxTypeTransfer   
)

func (t TransactionType) IsValid() bool {
	return t >= TxTypeDeposit && t <= TxTypeTransfer
}


func (t TransactionType) String() string {
	names := map[TransactionType]string{
		TxTypeDeposit:  "deposit",
		TxTypeWithdraw: "withdraw",
		TxTypeTransfer: "transfer",
	}
	return names[t]
}

type TransactionStatus uint

const (
	TxStatusPending   TransactionStatus = iota + 1 
	TxStatusCompleted                              
	TxStatusFailed                                 
	TxStatusCancelled     
)

func (s TransactionStatus) IsValid() bool {
	return s >= TxStatusPending && s <= TxStatusCancelled
}

func (s TransactionStatus) String() string {
	names := map[TransactionStatus]string{
		TxStatusPending:   "pending",
		TxStatusCompleted: "completed",
		TxStatusFailed:    "failed",
		TxStatusCancelled: "cancelled",
	}
	return names[s]
}

type Transaction struct {
	ID         uuid.UUID         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FromUserID *uuid.UUID        `json:"from_user_id" gorm:"type:uuid;index"`
	ToUserID   *uuid.UUID        `json:"to_user_id" gorm:"type:uuid;index"`
	Amount     float64           `json:"amount" gorm:"type:decimal(15,2);not null"`
	Type       TransactionType   `json:"type" gorm:"type:smallint;not null"`
	Status     TransactionStatus `json:"status" gorm:"type:smallint;default:1"`
	CreatedAt  time.Time         `json:"created_at"`

	FromUser *User `json:"from_user" gorm:"foreignKey:FromUserID"`
	ToUser   *User `json:"to_user" gorm:"foreignKey:ToUserID"`
}

func (Transaction) TableName() string {
	return "transactions"
}