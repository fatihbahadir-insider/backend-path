package transformer

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/constants"
)

func TransactionTransformer(tx *models.Transaction) dto.TransactionResponse {
	response := dto.TransactionResponse{
		ID:        tx.ID,
		Amount:    tx.Amount,
		Type:      tx.Type.String(),
		Status:    tx.Status.String(),
		CreatedAt: tx.CreatedAt.Format(constants.TimestampFormat),
	}

	if tx.FromUserID != nil {
		fromID := tx.FromUserID.String()
		response.FromUserID = &fromID
	}
	if tx.ToUserID != nil {
		toID := tx.ToUserID.String()
		response.ToUserID = &toID
	}

	return response
}

func TransactionListTransformer(transactions []models.Transaction) []dto.TransactionResponse {
	result := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		result[i] = TransactionTransformer(&tx)
	}
	return result
}