package services

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/app/repository"
	"backend-path/app/transformer"
	"backend-path/app/workers"
	"backend-path/utils"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ITransactionService interface {
	Credit(ctx *fiber.Ctx, req dto.CreditRequest, userID uuid.UUID) error
	Debit(ctx *fiber.Ctx, req dto.DebitRequest, userID uuid.UUID) error
	Transfer(ctx *fiber.Ctx, req dto.TransferRequest, fromUserID uuid.UUID) error
	GetByID(ctx *fiber.Ctx, id uuid.UUID, userID uuid.UUID) error
	GetHistory(ctx *fiber.Ctx, userID uuid.UUID) error
	GetStats(ctx *fiber.Ctx) error
}

type TransactionService struct {
	transactionRepo repository.ITransactionRepository
	balanceRepo     repository.IBalanceRepository
	auditRepo       repository.IAuditLogRepository
	workerPool      *workers.TransactionWorkerPool
}

var transactionServiceInstance *TransactionService

func NewTransactionService() *TransactionService {
	if transactionServiceInstance == nil {
		svc := &TransactionService{
			transactionRepo: repository.NewTransactionRepository(),
			balanceRepo: repository.NewBalanceRepository(),
			auditRepo: repository.NewAuditRepository(),
		}

		svc.workerPool = workers.NewTransactionWorkerPool(5, 100, svc.processTransaction)
		svc.workerPool.Start()

		transactionServiceInstance = svc
	}

	return transactionServiceInstance
}

func (s *TransactionService) processTransaction(job workers.TransactionJob) workers.TransactionResult {
	db := s.transactionRepo.GetDB()
	var resultTx *models.Transaction

	err := db.Transaction(func(tx *gorm.DB) error {
		transaction := &models.Transaction{
			ID: job.ID,
			FromUserID: job.FromUserID,
			ToUserID: job.ToUserID,
			Amount: job.Amount,
			Type: job.Type,
			Status: models.TxStatusPending,
		}

		if err := s.transactionRepo.Create(tx, transaction); err != nil {
			return err
		}

		var processErr error
		switch job.Type {
		case models.TxTypeDeposit:
			processErr = s.processDeposit(tx, transaction)
		case models.TxTypeWithdraw:
			processErr = s.processWithdraw(tx, transaction)
		case models.TxTypeTransfer:
			processErr = s.processTransfer(tx, transaction)
		}

		if processErr != nil {
			transaction.Fail()
			s.transactionRepo.Update(tx, transaction)
			return processErr
		}

		transaction.Complete()
		if err := s.transactionRepo.Update(tx, transaction); err != nil {
			return err
		}

		resultTx = transaction
		return nil
	})

	return workers.TransactionResult{
		Transaction: resultTx,
		Error: err,
	}
}
	
func (s *TransactionService) processDeposit(tx *gorm.DB, transaction *models.Transaction) error {
	balance, err := s.balanceRepo.FindByUserIDForUpdate(tx, *transaction.ToUserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	previousAmount := 0.0
	if balance != nil {
		previousAmount = balance.Amount
	}

	newBalance := &models.Balance{
		UserID: *transaction.ToUserID,
		Amount: previousAmount + transaction.Amount,
		LastUpdatedAt: time.Now(),
	}

	utils.Logger.Info("newBalance", zap.Float64("amount", newBalance.Amount))

	if err := s.balanceRepo.Upsert(tx, newBalance); err != nil {
		return err
	}

	s.logBalanceChange(transaction.ToUserID, models.ActionDeposit, previousAmount, newBalance.Amount, transaction.Amount, nil, &transaction.ID)

	return nil
}

func (s *TransactionService) processWithdraw(tx *gorm.DB, transaction *models.Transaction) error {
	balance, err := s.balanceRepo.FindByUserIDForUpdate(tx, *transaction.FromUserID)
	if err != nil {
		return err
	}

	previousAmount := balance.Amount
	balance.Amount -= transaction.Amount
	balance.LastUpdatedAt = time.Now()

	if balance.Amount < 0 {
		return errors.New("insufficient balance")
	}

	if err := s.balanceRepo.Update(tx, balance); err != nil {
		return err
	}

	s.logBalanceChange(transaction.FromUserID, models.ActionWithdraw, previousAmount, balance.Amount, transaction.Amount, nil, &transaction.ID)

	return nil
}

func (s *TransactionService) processTransfer(tx *gorm.DB, transaction *models.Transaction) error {
	fromID := *transaction.FromUserID
	toID := *transaction.ToUserID

	var firstID, secondID uuid.UUID
	if fromID.String() < toID.String() {
		firstID, secondID = fromID, toID
	} else {
		firstID, secondID = toID, fromID
	}
	
	firstBalance, err := s.balanceRepo.FindByUserIDForUpdate(tx, firstID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	secondBalance, err := s.balanceRepo.FindByUserIDForUpdate(tx, secondID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	
	var fromBalance, toBalance *models.Balance
	if firstID == fromID {
		fromBalance, toBalance = firstBalance, secondBalance
	} else {
		fromBalance, toBalance = secondBalance, firstBalance
	}

	if fromBalance == nil || fromBalance.Amount < transaction.Amount {
		return errors.New("insufficient balance")
	}

	fromPrevious := fromBalance.Amount
	fromBalance.Amount -= transaction.Amount
	fromBalance.LastUpdatedAt = time.Now()
	if err := s.balanceRepo.Update(tx, fromBalance); err != nil {
		return err
	}

	toPrevious := 0.0
	if toBalance != nil {
		toPrevious = toBalance.Amount
	}

	newToBalance := &models.Balance{
		UserID: *transaction.ToUserID,
		Amount: toPrevious + transaction.Amount,
		LastUpdatedAt: time.Now(),
	}

	if err := s.balanceRepo.Upsert(tx, newToBalance); err != nil {
		return err
	}

	s.logBalanceChange(&fromID, models.ActionTransferOut, fromPrevious, fromBalance.Amount, transaction.Amount, &toID, &transaction.ID)
	s.logBalanceChange(&toID, models.ActionTransferIn, toPrevious, newToBalance.Amount, transaction.Amount, &fromID, &transaction.ID)

	return nil
}

func (s *TransactionService) logBalanceChange(userID *uuid.UUID, action models.AuditAction, prev, new, change float64, relatedUserID, txID *uuid.UUID) {
	details := map[string]interface{}{
		"previous_amount": prev,
		"new_amount":      new,
		"change_amount":   change,
	}
	if relatedUserID != nil {
		details["related_user_id"] = relatedUserID.String()
	}
	if txID != nil {
		details["transaction_id"] = txID.String()
	}

	detailsJSON, _ := json.Marshal(details)
	go s.auditRepo.Create(&models.AuditLog{
		ID: uuid.New(),
		EntityType: models.EntityBalance,
		EntityID:   *userID,
		Action:     action,
		Details:    string(detailsJSON),
	})
}

func (s *TransactionService) Credit(ctx *fiber.Ctx, req dto.CreditRequest, userID uuid.UUID) error {
	if errs := utils.ValidateStruct(req); errs != nil {
		return utils.JsonErrorValidationFields(ctx, errs)
	}

	job := workers.TransactionJob{
		ID: uuid.New(),
		Type: models.TxTypeDeposit,
		ToUserID: &userID,
		Amount: req.Amount,
	}

	result := s.workerPool.SubmitAndWait(job)
	if result.Error != nil {
		return utils.JsonError(ctx, result.Error, "E_CREDIT_FAILED");
	}

	return utils.JsonSuccess(ctx, transformer.TransactionTransformer(result.Transaction))
}

func (s *TransactionService) Debit(ctx *fiber.Ctx, req dto.DebitRequest, userID uuid.UUID) error {
	if errs := utils.ValidateStruct(req); errs != nil {
		return utils.JsonErrorValidationFields(ctx, errs)
	}

	job := workers.TransactionJob{
		ID: uuid.New(),
		Type: models.TxTypeWithdraw,
		FromUserID: &userID,
		Amount: req.Amount,
	}

	result := s.workerPool.SubmitAndWait(job)
	if result.Error != nil {
		return utils.JsonError(ctx, result.Error, "E_DEBIT_FAILED");
	}

	return utils.JsonSuccess(ctx, transformer.TransactionTransformer(result.Transaction))
}

func (s *TransactionService) Transfer(ctx *fiber.Ctx, req dto.TransferRequest, fromUserID uuid.UUID) error {
	if errs := utils.ValidateStruct(req); errs != nil {
		return utils.JsonErrorValidationFields(ctx, errs)
	}

	if req.ToUserID ==  fromUserID {
		return utils.JsonError(ctx, errors.New("cannot transfer to yourself"), "E_TRANSFER_SELF");
	}

	job := workers.TransactionJob{
		ID: uuid.New(),
		Type: models.TxTypeTransfer,
		FromUserID: &fromUserID,
		ToUserID: &req.ToUserID,
		Amount: req.Amount,
	}

	result := s.workerPool.SubmitAndWait(job)
	if result.Error != nil {
		return utils.JsonError(ctx, result.Error, "E_TRANSFER_FAILED");
	}

	return utils.JsonSuccess(ctx, transformer.TransactionTransformer(result.Transaction))
}

func (s *TransactionService) GetByID(ctx *fiber.Ctx, id uuid.UUID, userID uuid.UUID) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TRANSACTION_NOT_FOUND")
	}

	if transaction.FromUserID != nil && *transaction.FromUserID != userID {
		return utils.JsonErrorUnauthorized(ctx, errors.New("unauthorized access"))
	}

	return utils.JsonSuccess(ctx, transformer.TransactionTransformer(transaction))
}

func (s *TransactionService) GetHistory(ctx *fiber.Ctx, userID uuid.UUID) error {
	pagination := utils.GetPagination(ctx)
	transactions, total, err := s.transactionRepo.FindByUserID(userID, pagination.Limit, pagination.GetOffset())
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TRANSACTION_HISTORY")
	}

	response := dto.NewPaginatedResponse(
		transformer.TransactionListTransformer(transactions),
		pagination.Page, pagination.Limit, total,
	)
	return utils.JsonSuccess(ctx, response)
}

func (s *TransactionService) GetStats(ctx *fiber.Ctx) error {
	stats := s.workerPool.GetStats()
	return utils.JsonSuccess(ctx, dto.TransactionStatsResponse{
		TotalProcessed:   stats.TotalProcessed,
		TotalSuccessful:  stats.TotalSuccessful,
		TotalFailed:      stats.TotalFailed,
		PendingInQueue:   s.workerPool.QueueLength(),
		TotalCredited:    float64(stats.TotalCredited) / 100,
		TotalDebited:     float64(stats.TotalDebited) / 100,
		TotalTransferred: float64(stats.TotalTransferred) / 100,
	})
}