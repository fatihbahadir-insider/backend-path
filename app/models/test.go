package models

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	GUID        uuid.UUID `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"type:varchar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Test) TableName() string {
	return "tests"
}