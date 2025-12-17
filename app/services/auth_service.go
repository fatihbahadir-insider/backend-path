package services

import (
	"backend-path/app/dto"
	"backend-path/app/middlewares"
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/constants"
	"backend-path/utils"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {}

func (s *AuthService) userRepo() *repository.UserRepository {
	return new(repository.UserRepository)
}

func (s *AuthService) Authenticate(ctx *fiber.Ctx, req dto.LoginRequest) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}

	user, err := s.userRepo().FindByEmail(req.Email);
	if err != nil {
		err := errors.New("username or password is wrong")
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		err := errors.New("username or password is wrong")
		return utils.JsonErrorValidation(ctx, err)
	}

	expireHour, _ := time.ParseDuration(os.Getenv("JWT_EXPIRES") + "h")
	expiresAt := time.Now().Add(time.Hour * expireHour).Unix()
	token, err := s.generateToken(user.ID.String(), expiresAt)

	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TOKEN_GENERATE")
	}


	return utils.JsonSuccess(ctx, dto.AuthResponse{
		Token: token,
	})
}

func (s *AuthService) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}

	userRepo := s.userRepo()
	if userRepo.IsExist(req.Email) {
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

	if err := userRepo.Insert(user); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_CREATE")
	}

	return utils.JsonSuccess(ctx, user)
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *AuthService) generateToken(userGUID string, expiresAt int64) (string, error) {
	claims := middlewares.JwtCustomClaims{
		Issuer: userGUID,
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
