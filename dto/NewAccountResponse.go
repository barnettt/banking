package dto

type NewAccountResponse struct {
	AccountId string `json:"account_id"  xml:"account_id"`
}

func GetAccountResponse(inAccountId string) NewAccountResponse {
	return NewAccountResponse{AccountId: inAccountId}
}
