package dto

import "github.com/google/uuid"

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"omitempty,email,max=100"`
	Role   	 string `json:"role" validate:"omitempty,oneof=user admin mod"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total uint64 `json:"total"`
}