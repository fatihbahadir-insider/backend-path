package transformer

import (
	"backend-path/app/dto"
	"backend-path/app/models"
	"backend-path/constants"
	"encoding/json"
)

func BalanceTransformer(balance *models.Balance) dto.BalanceResponse {
	return dto.BalanceResponse{
		UserID:        balance.UserID,
		Amount:        balance.Amount,
		LastUpdatedAt: balance.LastUpdatedAt.Format(constants.TimestampFormat),
	}
}

func BalanceHistoryTransformer(logs []models.AuditLog) dto.BalanceHistoryResponse {
	history := make([]dto.BalanceHistoryItem, 0, len(logs))

	for _, log := range logs {
		var details map[string]interface{}
		json.Unmarshal([]byte(log.Details), &details)

		item := dto.BalanceHistoryItem{
			Action:    log.Action.String(),
			CreatedAt: log.CreatedAt.Format(constants.TimestampFormat),
		}

		if v, ok := details["previous_amount"].(float64); ok {
			item.PreviousAmount = v
		}
		if v, ok := details["new_amount"].(float64); ok {
			item.NewAmount = v
		}
		if v, ok := details["change_amount"].(float64); ok {
			item.ChangeAmount = v
		}
		if v, ok := details["related_user_id"].(string); ok {
			item.RelatedUserID = &v
		}
		if v, ok := details["transaction_id"].(string); ok {
			item.TransactionID = &v
		}

		history = append(history, item)
	}

	return dto.BalanceHistoryResponse{
		History: history,
	}
}