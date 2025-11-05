package metrics

import (
	"encoding/json"
	"fmt"
	"lib/models"
)

// ParseBlockJSON парсит JSON блока в models.Block без квитанций
func ParseBlockJSON(blockRaw json.RawMessage) models.BlockDTO {
	var rpcDTO models.BlockRPCDTO
	if err := json.Unmarshal(blockRaw, &rpcDTO); err != nil {
		fmt.Printf("failed to unmarshal block JSON: %v\n", err)
		return models.BlockDTO{}
	}
	return ConvertBlockDTO(rpcDTO)
}

func ParseBlockReceiptsJSON(receiptsRaw json.RawMessage) []models.ReceiptDTO {
	var rpcDTOs []models.ReceiptRpcDTO

	if err := json.Unmarshal(receiptsRaw, &rpcDTOs); err != nil {
		fmt.Printf("failed to unmarshal block receipts to DTO: %v\n", err)
		return nil
	}

	receipts := make([]models.ReceiptDTO, 0, len(rpcDTOs))
	for _, rpcDTO := range rpcDTOs {
		receipts = append(receipts, ConvertReceiptDTO(rpcDTO))
	}
	return receipts
}
