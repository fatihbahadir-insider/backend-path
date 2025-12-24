package services

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/app/transformer"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IUserService interface {
	GetAll(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx, id uuid.UUID) error
	Update(ctx *fiber.Ctx, id uuid.UUID, req dto.UpdateUserRequest) error
	Delete(ctx *fiber.Ctx, id uuid.UUID) error
}

type UserService struct {
	userRepo repository.IUserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}


func (s *UserService) GetAll(ctx *fiber.Ctx) error {
	pagination := utils.GetPagination(ctx)

	users, total, err := s.userRepo.FindAll(pagination.Limit, pagination.GetOffset())
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_LIST")
	}

	response := dto.NewPaginatedResponse(
		transformer.UserListTransformer(users),
		pagination.Page,
		pagination.Limit,
		total,
	)

	return utils.JsonSuccess(ctx, response)
}

func (s *UserService) GetByID(ctx *fiber.Ctx, id uuid.UUID) error {
	user ,err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("user not found"))
	}

	return utils.JsonSuccess(ctx, transformer.UserTransformer(user))
}

func (s *UserService) Update(ctx *fiber.Ctx, id uuid.UUID, req dto.UpdateUserRequest) error {
	if errs := utils.ValidateStruct(req); errs != nil {
		return utils.JsonErrorValidationFields(ctx, errs)
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("user not found"))
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Role != "" {
		user.RoleID = models.ToRole(req.Role)
	}

	if err := s.userRepo.Update(user); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_UPDATE")
	}

	return utils.JsonSuccess(ctx, transformer.UserTransformer(user))
}

func (s *UserService) Delete(ctx *fiber.Ctx, id uuid.UUID) error {
	currentUserID := ctx.Locals("user_auth").(string)
	if currentUserID == id.String() {
		return utils.JsonErrorForbidden(ctx, errors.New("cannot delete yourself"))
	}

	if _, err := s.userRepo.FindByID(id); err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("user not found"))
	}

	if err := s.userRepo.Delete(id); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_DELETE")
	}

	return utils.JsonSuccess(ctx, fiber.Map{"message": "user deleted"})
}