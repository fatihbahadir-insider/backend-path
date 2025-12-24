package services

import (
	"backend-path/app/dto"
	"backend-path/app/repository"
	"backend-path/app/transformer"
	"backend-path/constants"
	"backend-path/utils"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IBalanceService interface {
	GetCurrent(ctx *fiber.Ctx, userID uuid.UUID) error
	GetHistorical(ctx *fiber.Ctx, userID uuid.UUID) error
	GetAtTime(ctx *fiber.Ctx, req dto.BalanceAtTimeRequest, userID uuid.UUID) error
}

type BalanceService struct {
	balanceRepo repository.IBalanceRepository
}

func NewBalanceService() *BalanceService {
	return &BalanceService{
		balanceRepo: repository.NewBalanceRepository(),
	}
}

func (s *BalanceService) GetCurrent(ctx *fiber.Ctx, userID uuid.UUID) error {
	balance, err := s.balanceRepo.FindByUserID(userID)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	return utils.JsonSuccess(ctx, transformer.BalanceTransformer(balance))
}

func (s *BalanceService) GetHistorical(ctx *fiber.Ctx, userID uuid.UUID) error {
	pagination := utils.GetPagination(ctx)

	logs, total, err := s.balanceRepo.GetBalanceHistory(
		userID,
		pagination.Limit,
		pagination.GetOffset(),
	)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	response := dto.NewPaginatedResponse(
		transformer.BalanceHistoryTransformer(logs),
		pagination.Page,
		pagination.Limit,
		total,
	)

	return utils.JsonSuccess(ctx, response)
}

func (s *BalanceService) GetAtTime(ctx *fiber.Ctx, req dto.BalanceAtTimeRequest, userID uuid.UUID) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}
	
	timestamp := time.Unix(req.Timestamp, 0)


	log, err := s.balanceRepo.GetBalanceAtTime(userID, timestamp)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	var details map[string]interface{}
	json.Unmarshal([]byte(log.Details), &details)

	amount := 0.0
	if v, ok := details["new_amount"].(float64); ok {
		amount = v
	}

	return utils.JsonSuccess(ctx, dto.BalanceAtTimeResponse{
		UserID: userID,
		Amount: amount,
		AsOf:    log.CreatedAt.Format(constants.TimestampFormat),
		IsExact: log.CreatedAt.Equal(timestamp),
	})
}