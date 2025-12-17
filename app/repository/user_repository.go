package repository

import "backend-path/app/models"

type UserRepository struct {}

func (r *UserRepository) Insert(user *models.User) error {
	if err := DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IsExist(email string) bool {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}

	return true
}