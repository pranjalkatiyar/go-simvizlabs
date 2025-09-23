package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"simvizlab-backend/models"
	"simvizlab-backend/utils"

	"simvizlab-backend/services"

	"github.com/gin-gonic/gin"
)

type TransactionRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

func GetTransactionInfo(ctx *gin.Context) {
	var req TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid transaction_id"})
		return
	}

	jwtToken, err := utils.GenerateAppStoreJWT()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	data, err := services.FetchTransaction(jwtToken, req.TransactionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	var response models.TransactionInfoResponse

	if err := json.Unmarshal(data, &response); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse transaction response"})
		return
	}

	results, err := utils.DecodeSignedTransactionInfo(string(response.SignedTransactionInfo))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction data"})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func GetHistoryInfo(ctx *gin.Context) {
	var req TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid transaction_id"})
		return
	}

	jwtToken, err := utils.GenerateAppStoreJWT()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	data, err := services.FetchTransactionHistory(jwtToken, req.TransactionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	var historyResponse models.HistoryResponse
	var AppleAppId int64
	if err := json.Unmarshal(data, &historyResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse history response"})
		return
	}
	AppleAppId = historyResponse.AppAppleId
	fmt.Println("AppleAppId:", AppleAppId)
	results, err := utils.DecodeSignedTransactionInfo(string(historyResponse.SignedTransactions[0]))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction data"})
		return
	}
	// ctx.Data(http.StatusOK, "application/json", data)

	ctx.JSON(http.StatusOK, results)
}
