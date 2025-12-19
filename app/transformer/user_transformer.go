package transformer

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/constants"
)

func UserTransformer(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.RoleID.String(),
		CreatedAt: user.CreatedAt.Format(constants.TimestampFormat),
	}
}

func UserListTransformer(users []models.User) []dto.UserResponse {
	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, UserTransformer(&user))
	}
	return response
}