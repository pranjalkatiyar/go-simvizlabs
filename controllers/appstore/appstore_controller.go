package controller

import (
	"net/http"

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

	jwtToken, err := utils.GenerateJWT()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	data, err := services.FetchTransaction(jwtToken, req.TransactionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	ctx.Data(http.StatusOK, "application/json", data)
}
