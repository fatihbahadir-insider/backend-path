package models

import (
	"time"

	"github.com/google/uuid"
)

type EntityType uint

const (
	EntityUser        EntityType = iota + 1 
	EntityTransaction                            
	EntityBalance                                  
	EntityRole                                     
)

func (e EntityType) IsValid() bool {
	return e >= EntityUser && e <= EntityRole
}

func (e EntityType) String() string {
	names := map[EntityType]string{
		EntityUser:        "user",
		EntityTransaction: "transaction",
		EntityBalance:     "balance",
		EntityRole:        "role",
	}
	return names[e]
}

type AuditAction uint

const (
	ActionCreate AuditAction = iota + 1 
	ActionUpdate                            
	ActionDelete                            
	ActionLogin   
	ActionRegister                          
	ActionLogout             
	ActionRefreshToken      
	ActionDeposit
	ActionWithdraw
	ActionTransferIn
	ActionTransferOut                      
)

func (a AuditAction) IsValid() bool {
	return a >= ActionCreate && a <= ActionLogout
}

func (a AuditAction) String() string {
	names := map[AuditAction]string{
		ActionCreate: "create",
		ActionUpdate: "update",
		ActionDelete: "delete",
		ActionLogin:  "login",
		ActionRegister: "register",
		ActionLogout: "logout",
		ActionRefreshToken: "refresh_token",
	}
	return names[a]
}

type AuditLog struct {
	ID         uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EntityType EntityType  `json:"entity_type" gorm:"type:smallint;not null;index"`
	EntityID   uuid.UUID   `json:"entity_id" gorm:"type:uuid;index"`
	Action     AuditAction `json:"action" gorm:"type:smallint;not null"`
	Details    string      `json:"details" gorm:"type:jsonb"`
	CreatedAt  time.Time   `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}