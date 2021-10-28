package dto

type TransactionResponse struct {
	TransactionId     string  `json:"transaction_id" xml:"transaction_id"`
	TransactionAmount float64 `json:"transaction_amount" xml:"transaction_amount"`
	TransactionType   string  `json:"transaction_type" xml:"transaction_type"`
	TransactionDate   string  `json:"transaction_date" xml:"transaction_date"`
	Balance           float64 `json:"balance" xml:"balance"`
}

func GetTransactionResponse(transactionId string) TransactionResponse {
	return TransactionResponse{TransactionId: transactionId}
}
