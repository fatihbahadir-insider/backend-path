package services

import (
	"backend-path/app/dto"
	"backend-path/app/repository"
	"backend-path/app/transformer"
	"backend-path/configs"
	"backend-path/constants"
	"backend-path/utils"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
)

type IBalanceService interface {
	GetCurrent(ctx *fiber.Ctx, userID uuid.UUID) error
	GetHistorical(ctx *fiber.Ctx, userID uuid.UUID) error
	GetAtTime(ctx *fiber.Ctx, req dto.BalanceAtTimeRequest, userID uuid.UUID) error
}

type BalanceService struct {
	balanceRepo repository.IBalanceRepository
	redisStorage *redis.Storage
}

func NewBalanceService() *BalanceService {
	return &BalanceService{
		balanceRepo: repository.NewBalanceRepository(),
		redisStorage: configs.RedisStorage,
	}
}

func (s *BalanceService) GetCurrent(ctx *fiber.Ctx, userID uuid.UUID) error {
	cacheKey := s.keyBalanceCurrentCache(userID)
	cacheData := s.getBalanceCurrentCache(cacheKey)

	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

	balance, err := s.balanceRepo.FindByUserID(userID)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	response := transformer.BalanceTransformer(balance)
	s.setCache(cacheKey, response)


	return utils.JsonSuccess(ctx, response)
}

func (s *BalanceService) GetHistorical(ctx *fiber.Ctx, userID uuid.UUID) error {
	pagination := utils.GetPagination(ctx)
	cacheKey := s.keyBalanceHistoryCache(userID, pagination.Page, pagination.Limit)
	cacheData := s.getBalanceHistoryCache(cacheKey)

	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
	}

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

	s.setCache(cacheKey, response)
	return utils.JsonSuccess(ctx, response)
}

func (s *BalanceService) GetAtTime(ctx *fiber.Ctx, req dto.BalanceAtTimeRequest, userID uuid.UUID) error {
	if errors := utils.ValidateStruct(req); errors != nil {
		return utils.JsonErrorValidationFields(ctx, errors)
	}
	
	cacheKey := s.keyBalanceAtTimeCache(userID, req.Timestamp)
	cacheData := s.getBalanceAtTimeCache(cacheKey)

	if cacheData != nil {
		return utils.JsonSuccess(ctx, cacheData)
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

	response := dto.BalanceAtTimeResponse{
		UserID: userID,
		Amount: amount,
		AsOf:    log.CreatedAt.Format(constants.TimestampFormat),
		IsExact: log.CreatedAt.Equal(timestamp),
	}

	s.setCache(cacheKey, response)
	return utils.JsonSuccess(ctx, response)
}

func (s *BalanceService) keyBalanceCurrentCache(userID uuid.UUID) string {
	return constants.CacheBalanceCurrent + "_" + userID.String()
}

func (s *BalanceService) keyBalanceHistoryCache(userID uuid.UUID, page, limit int) string {
	return constants.CacheBalanceHistory + "_" + userID.String() + "_" +
		strconv.Itoa(page) + "_" + strconv.Itoa(limit)
}

func (s *BalanceService) keyBalanceAtTimeCache(userID uuid.UUID, timestamp int64) string {
	return constants.CacheBalanceAtTime + "_" + userID.String() + "_" +
		strconv.FormatInt(timestamp, 10)
}

func (s *BalanceService) getBalanceCurrentCache(key string) *dto.BalanceResponse {
	if s.redisStorage == nil {
		return nil
	}

	data, err := s.redisStorage.Get(key)
	if err != nil || len(data) == 0 {
		utils.Logger.Info("NO CACHE BALANCE CURRENT FOUND FOR KEY " + key)
		return nil
	}

	var response dto.BalanceResponse
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Logger.Error("Error unmarshaling balance current cache: " + err.Error())
		return nil
	}

	utils.Logger.Info("GET CACHE BALANCE CURRENT FROM KEY " + key)
	return &response
}

func (s *BalanceService) getBalanceHistoryCache(key string) *dto.PaginatedResponse[dto.BalanceHistoryItem] {
	if s.redisStorage == nil {
		return nil
	}

	data, err := s.redisStorage.Get(key)
	if err != nil || len(data) == 0 {
		utils.Logger.Info("NO CACHE BALANCE HISTORY FOUND FOR KEY " + key)
		return nil
	}

	var response dto.PaginatedResponse[dto.BalanceHistoryItem]
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Logger.Error("Error unmarshaling balance history cache: " + err.Error())
		return nil
	}

	utils.Logger.Info("GET CACHE BALANCE HISTORY FROM KEY " + key)
	return &response
}

func (s *BalanceService) getBalanceAtTimeCache(key string) *dto.BalanceAtTimeResponse {
	if s.redisStorage == nil {
		return nil
	}

	data, err := s.redisStorage.Get(key)
	if err != nil || len(data) == 0 {
		utils.Logger.Info("NO CACHE BALANCE AT TIME FOUND FOR KEY " + key)
		return nil
	}

	var response dto.BalanceAtTimeResponse
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Logger.Error("Error unmarshaling balance at time cache: " + err.Error())
		return nil
	}

	utils.Logger.Info("GET CACHE BALANCE AT TIME FROM KEY " + key)
	return &response
}

func (s *BalanceService) setCache(key string, data interface{}) {
	if s.redisStorage == nil {
		utils.Logger.Error("REDIS STORAGE IS NULL")
		return
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		utils.Logger.Error("Error marshaling user list cache: " + err.Error())
		return
	}

	if err := s.redisStorage.Set(key, dataJSON, 12 * time.Hour); err != nil {
		utils.Logger.Error("REDIS KEY " + key + " ERROR: " + err.Error())
	}

	utils.Logger.Info("SET CACHE BALANCE LIST TO KEY " + key)
}