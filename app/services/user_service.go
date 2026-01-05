package services

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/app/transformer"
	"backend-path/configs"
	"backend-path/constants"
	"backend-path/utils"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
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
	redisStorage *redis.Storage
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
		redisStorage: configs.RedisStorage,
	}
}


func (s *UserService) GetAll(ctx *fiber.Ctx) error {
	pagination := utils.GetPagination(ctx)

	cacheKey := s.keyUserListCache(pagination.Page, pagination.Limit)
	cacheData := s.getUserListCache(cacheKey)
	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

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

	if response.Data != nil {
		s.setCache(cacheKey, response)
	}

	return utils.JsonSuccess(ctx, response)
}

func (s *UserService) GetByID(ctx *fiber.Ctx, id uuid.UUID) error {
	cacheKey := s.keyUserDetailCache(id.String())
	cacheData := s.getUserDetailCache(cacheKey)

	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

	user ,err := s.userRepo.FindByID(id)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("user not found"))
	}

	response := transformer.UserTransformer(user)

	if response.ID != uuid.Nil {
		s.setCache(cacheKey, response)
	}

	return utils.JsonSuccess(ctx, response)
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

	userDetail := transformer.UserTransformer(user)
	cacheKey := s.keyUserDetailCache(id.String())
	s.setCache(cacheKey, userDetail)

	s.resetUserListCache()

	return utils.JsonSuccess(ctx, userDetail)
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


	cacheKey := s.keyUserDetailCache(id.String())
	s.redisStorage.Delete(cacheKey)

	s.resetUserListCache()
	
	return utils.JsonSuccess(ctx, fiber.Map{"message": "user deleted"})
}

func (s *UserService) keyUserListCache(page int, limit int) string {
	return constants.CacheUserList + "_" + strconv.Itoa(page) + "_" + strconv.Itoa(limit)
}

func (s *UserService) keyUserDetailCache(uuid string) string {
	return constants.CacheUserDetail + "_" + uuid
}

func (s *UserService) getUserListCache(key string) *dto.PaginatedResponse[dto.UserResponse] {
	cacheUsers, err := s.redisStorage.Get(key)
	if err != nil || len(cacheUsers) == 0 {
		utils.Logger.Info("❌ NO CACHE USER LIST FOUND FOR KEY " + key)
		return nil
	}

	var response dto.PaginatedResponse[dto.UserResponse]
	if err := json.Unmarshal(cacheUsers, &response); err != nil {
		utils.Logger.Error("Error unmarshaling user list cache: " + err.Error())
		return nil
	}

	utils.Logger.Info("✅ GET CACHE USER LIST FROM KEY " + key)
	return &response
}

func (s *UserService) getUserDetailCache(key string) (userDetail *dto.UserResponse) {
	cacheUser, _ := s.redisStorage.Get(key)
	if len(cacheUser) > 0 {
		utils.Logger.Info("✅ GET CACHE USER DETAIL FROM KEY " + key)
		if err := json.Unmarshal(cacheUser, &userDetail); err != nil {
			return nil
		}

		return userDetail
	}

	return nil
}

func (s *UserService) setCache(key string, data interface{}) {
	if s.redisStorage == nil {
		utils.Logger.Error("❌ REDIS STORAGE IS NULL")
		return
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		utils.Logger.Error("Error marshaling user list cache: " + err.Error())
		return
	}

	if err := s.redisStorage.Set(key, dataJSON, 12 * time.Hour); err != nil {
		utils.Logger.Error("❌ REDIS KEY " + key + " ERROR: " + err.Error())
	}

	utils.Logger.Info("✅ SET CACHE USER LIST TO KEY " + key)
}

func (s *UserService) resetUserListCache() {
	if s.redisStorage == nil {
		return
	}

	commonPages := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30, 40, 50}
	commonLimits := []int{10, 20, 25, 50, 100}
	
	deletedCount := 0
	for _, page := range commonPages {
		for _, limit := range commonLimits {
			key := s.keyUserListCache(page, limit)
			if err := s.redisStorage.Delete(key); err == nil {
				deletedCount++
			}
		}
	}
	
	utils.Logger.Info("✅ RESET CACHE USER LIST - Deleted " + strconv.Itoa(deletedCount) + " cache keys")
}