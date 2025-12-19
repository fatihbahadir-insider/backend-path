package repository

import (
	"backend-path/app/models"

	"github.com/google/uuid"
)

type IUserRepository interface {
	Insert(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindAll() ([]models.User, uint64, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	IsExist(email string) bool
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindAll() ([]models.User, uint64, error) {
	var users []models.User
	var total int64

	DB.Model(&models.User{}).Count(&total)
	if err := DB.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, uint64(total), nil
}

func (r *UserRepository) Insert(user *models.User) error {
	if err := DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	if err := DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
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

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
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