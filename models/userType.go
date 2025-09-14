package models

type CreateUserRequest struct {
	AppleAppId            int64  `json:"appleAppIid" binding:"required"`
	TransactionId         string `json:"transactionId""`
	OriginalTransactionId string `json:"originalTransactionId" binding:"required"`
}
