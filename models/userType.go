package models

type CreateUserRequest struct {
	AppleAppId    int64  `json:"apple_app_id" binding:"required"`
	TransactionId string `json:"transaction_id" binding:"required"`
}
