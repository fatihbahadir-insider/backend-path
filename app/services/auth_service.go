package services

import (
	"backend-path/app/dto"
	"backend-path/app/middlewares"
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/constants"
	"backend-path/utils"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Authenticate(ctx *fiber.Ctx, req dto.LoginRequest) error
	Register(ctx *fiber.Ctx, req dto.RegisterRequest) error
	RefreshToken(ctx *fiber.Ctx) error
}

type AuthService struct {
	userRepo repository.IUserRepository
	auditRepo repository.IAuditLogRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:  repository.NewUserRepository(),
		auditRepo: repository.NewAuditRepository(),
	}
}

func (s *AuthService) Authenticate(ctx *fiber.Ctx, req dto.LoginRequest) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}

	user, err := s.userRepo.FindByEmail(req.Email);
	if err != nil {
		s.logAuth(ctx, nil, models.ActionLogin, map[string]interface{}{
			"email": req.Email,
			"status": "failed",
			"reason": "user not found",
		})
		err := errors.New("username or password is wrong")
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logAuth(ctx, &user.ID, models.ActionLogin, map[string]interface{}{
			"email":  req.Email,
			"status": "failed",
			"reason": "wrong password",
		})
		err := errors.New("username or password is wrong")
		return utils.JsonErrorValidation(ctx, err)
	}


	token, err := s.generateToken(user.ID.String(), user.RoleID)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TOKEN_GENERATE")
	}

	s.logAuth(ctx, &user.ID, models.ActionLogin, map[string]interface{}{
		"email":  user.Email,
		"status": "success",
	})
	return utils.JsonSuccess(ctx, dto.AuthResponse{
		Token: token,
	})
}

func (s *AuthService) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}

	if s.userRepo.IsExist(req.Email) {
		s.logAuth(ctx, nil, models.ActionRegister, map[string]interface{}{
			"email": req.Email,
			"status": "failed",
			"reason": "email already exists",
		})
		return utils.JsonErrorValidation(ctx, constants.ErrEmailExist)
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_HASH_PASSWORD")
	}

	user := &models.User{
		Username: req.Username,
		Email: req.Email,
		PasswordHash: hashedPassword,
		RoleID: models.RoleUser,
	}

	if err := s.userRepo.Insert(user); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_CREATE")
	}

	s.logAuth(ctx, &user.ID, models.ActionRegister, map[string]interface{}{
		"email": user.Email,
		"status": "success",
	})
	return utils.JsonSuccess(ctx, user)
}

func (s *AuthService) RefreshToken(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("user not found"))
	}

	token, err := s.generateToken(user.ID.String(), user.RoleID)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TOKEN_GENERATE")
	}

	s.logAuth(ctx, &user.ID, models.ActionRefreshToken, map[string]interface{}{
		"email": user.Email,
		"status": "success",
	})

	return utils.JsonSuccess(ctx, dto.AuthResponse{
		Token: token,
	})
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *AuthService) generateToken(userUUID string, role models.Role) (string, error) {
	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES"))
	if expireHours == 0 {
		expireHours = 24
	}
	expiresAt := time.Now().Add(time.Duration(expireHours) * time.Hour).Unix()

	claims := middlewares.JwtCustomClaims{
		Issuer: userUUID,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (s *AuthService) logAuth(ctx *fiber.Ctx, userID *uuid.UUID, action models.AuditAction, details map[string]interface{}) {
	details["ip"] = ctx.IP()
	details["user_agent"] = string(ctx.Request().Header.UserAgent())

	detailsJSON, _ := json.Marshal(details)

	log := &models.AuditLog{
		EntityType: models.EntityUser,
		Action: action,
		Details: string(detailsJSON),
	}

	if userID != nil {
		log.EntityID = *userID
	}
	
	go s.auditRepo.Create(log)
}