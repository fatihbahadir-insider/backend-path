package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username     string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	RoleID       Role      `json:"role_id" gorm:"type:smallint;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Role         Role        `json:"role" gorm:"foreignKey:RoleID"`
}

func (User) TableName() string {
	return "users"
}