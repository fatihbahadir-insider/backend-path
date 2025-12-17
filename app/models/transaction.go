package models

import (
	"errors"
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

var validTransitions = map[TransactionStatus][]TransactionStatus{
	TxStatusPending: {TxStatusCompleted, TxStatusFailed, TxStatusCancelled},
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

func (t *Transaction) CanTransitionTo(newStatus TransactionStatus) bool {
	allowed, exists := validTransitions[t.Status]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

func (t *Transaction) Complete() error {
	if !t.CanTransitionTo(TxStatusCompleted) {
		return errors.New("cannot complete: transaction is not pending")
	}
	t.Status = TxStatusCompleted
	return nil
}

func (t *Transaction) Fail() error {
	if !t.CanTransitionTo(TxStatusFailed) {
		return errors.New("cannot fail: transaction is not pending")
	}
	t.Status = TxStatusFailed
	return nil
}

func (t *Transaction) Cancel() error {
	if !t.CanTransitionTo(TxStatusCancelled) {
		return errors.New("cannot cancel: transaction is not pending")
	}
	t.Status = TxStatusCancelled
	return nil
}

func (t *Transaction) IsPending() bool {
	return t.Status == TxStatusPending
}

func (t *Transaction) IsFinal() bool {
	return t.Status != TxStatusPending
}

func (t *Transaction) IsDeposit() bool {
	return t.Type == TxTypeDeposit
}

func (t *Transaction) IsWithdraw() bool {
	return t.Type == TxTypeWithdraw
}

func (t *Transaction) IsTransfer() bool {
	return t.Type == TxTypeTransfer
}