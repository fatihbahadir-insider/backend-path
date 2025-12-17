package dto

import "github.com/google/uuid"

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"omitempty,email,max=100"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}