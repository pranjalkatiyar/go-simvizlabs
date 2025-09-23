package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"simvizlab-backend/models"
	mongoRepo "simvizlab-backend/repository/mongo"
	"simvizlab-backend/services"
	"simvizlab-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers(ctx *gin.Context) {
	var users []*models.User
	err := mongoRepo.Get("users", &users)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No users found", "users": []models.User{}})
		return
	}

	response := gin.H{
		"users": users,
		"count": len(users),
	}

	ctx.JSON(http.StatusOK, response)
}

func GetOneUser(ctx *gin.Context) {
	var user models.User
	err := mongoRepo.GetOne("users", nil, &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func respondWithError(ctx *gin.Context, status int, message string, details ...string) {
	resp := gin.H{"error": message}
	if len(details) > 0 {
		resp["details"] = details[0]
	}
	ctx.JSON(status, resp)
}

func CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if req.AppleAppId <= 0 {
		respondWithError(ctx, http.StatusBadRequest, "AppleAppId must be a positive integer")
		return
	}

	var existingUser models.User
	err := mongoRepo.GetOne("users", bson.M{"appleAppId": req.AppleAppId}, &existingUser)

	if err == nil {
		respondWithError(ctx, http.StatusConflict, "User with this AppleAppId already exists")
		return
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		respondWithError(ctx, http.StatusInternalServerError, "Database error while checking user existence", err.Error())
		return
	}

	jwtToken, err := utils.GenerateAppStoreJWT()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to generate JWT", err.Error())
		return
	}

	data, err := services.FetchTransaction(jwtToken, req.OriginalTransactionId)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to fetch transaction", err.Error())
		return
	}

	var response models.TransactionInfoResponse
	if err := json.Unmarshal(data, &response); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to parse transaction response", err.Error())
		return
	}

	results, err := utils.DecodeSignedTransactionInfo(string(response.SignedTransactionInfo))
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to decode transaction data", err.Error())
		return
	}

	var decodedInfo models.JWSTransaction

	jsonBytes, err := json.Marshal(results)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to re-encode transaction map", err.Error())
		return
	}

	if err := json.Unmarshal(jsonBytes, &decodedInfo); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to parse decoded transaction data", err.Error())
		return
	}

	user := models.User{
		AppleAppId:            req.AppleAppId,
		OriginalTransactionId: req.OriginalTransactionId,
	}
	if err := mongoRepo.Save("users", &user); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to save user", err.Error())
		return
	}

	transaction := models.JWSTransaction{
		AppleAppId:            req.AppleAppId,
		AppTransactionId:      decodedInfo.AppTransactionId,
		OriginalTransactionId: req.OriginalTransactionId,
		WebOrderLineItemId:    decodedInfo.WebOrderLineItemId,
		BundleID:              decodedInfo.BundleID,
		ProductID:             decodedInfo.ProductID,
		PurchaseDate:          decodedInfo.PurchaseDate,
		ExpiresDate:           decodedInfo.ExpiresDate,
		Currency:              decodedInfo.Currency,
		OriginalPurchaseDate:  decodedInfo.OriginalPurchaseDate,
		Storefront:            decodedInfo.Storefront,
		Environment:           decodedInfo.Environment,
		SignedDate:            decodedInfo.SignedDate,
		Price:                 decodedInfo.Price,
		// Add more fields from results as needed
	}

	if err := mongoRepo.Save("transactions", &transaction); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to save transaction", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":        user,
		"transaction": transaction,
	})
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := mongoRepo.Update("users", map[string]interface{}{"_id": user.ID}, &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func CheckUserSubscriptionStatus(ctx *gin.Context) {
	appleIDStr := ctx.Query("appleAppId")
	if appleIDStr == "" {
		respondWithError(ctx, http.StatusBadRequest, "appleAppId is required")
		return
	}
	appleAppId, err := strconv.ParseInt(appleIDStr, 10, 64)
	if err != nil || appleAppId <= 0 {
		respondWithError(ctx, http.StatusBadRequest, "appleAppId must be a positive integer")
		return
	}

	var user models.User
	if err := mongoRepo.GetOne("users", bson.M{"appleAppId": appleAppId}, &user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondWithError(ctx, http.StatusNotFound, "user not found")
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "database error", err.Error())
		return
	}

	// Prefer explicit transactionId; otherwise use originalTransactionId (query or stored)
	transactionId := ctx.Query("transactionId")
	originalTransactionId := ctx.Query("originalTransactionId")
	if originalTransactionId == "" {
		originalTransactionId = user.OriginalTransactionId
	}

	jwtToken, err := utils.GenerateAppStoreJWT()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "failed to generate JWT", err.Error())
		return
	}

	var data []byte
	if transactionId != "" {
		data, err = services.FetchTransaction(jwtToken, transactionId)
	} else {
		if originalTransactionId == "" {
			respondWithError(ctx, http.StatusBadRequest, "originalTransactionId is required when transactionId is not provided")
			return
		}
		data, err = services.FetchTransactionHistory(jwtToken, originalTransactionId)
	}
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, "failed to fetch App Store data", err.Error())
		return
	}

	var jws string
	if transactionId != "" {
		var tr models.TransactionInfoResponse
		if err := json.Unmarshal(data, &tr); err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "failed to parse transaction response", err.Error())
			return
		}
		jws = tr.SignedTransactionInfo
	} else {
		var hist models.HistoryResponse
		if err := json.Unmarshal(data, &hist); err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "failed to parse history response", err.Error())
			return
		}
		if len(hist.SignedTransactions) == 0 {
			respondWithError(ctx, http.StatusNotFound, "no transactions found for originalTransactionId")
			return
		}
		jws = hist.SignedTransactions[0]
	}

	decoded, err := utils.DecodeSignedTransactionInfo(jws)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "failed to decode transaction JWS", err.Error())
		return
	}

	var tx models.JWSTransaction
	raw, _ := json.Marshal(decoded)
	if err := json.Unmarshal(raw, &tx); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "failed to parse decoded transaction", err.Error())
		return
	}

	// Determine active status by expiry (tolerate seconds or milliseconds)
	expires := tx.ExpiresDate
	nowMs := time.Now().UnixMilli()
	if expires > 0 && expires < 1_000_000_000_000 {
		expires *= 1000
	}
	active := expires == 0 || expires > nowMs

	ctx.JSON(http.StatusOK, gin.H{
		"exists":      true,
		"active":      active,
		"expiresAtMs": expires,
		"nowMs":       nowMs,
		"user":        user,
		"transaction": tx,
	})
}

func LoginAndCheckStatus(ctx *gin.Context) {
	type reqBody struct {
		AppleAppId            int64  `json:"appleAppId" binding:"required"`
		OriginalTransactionId string `json:"originalTransactionId" binding:"required"`
	}
	var req reqBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Find existing transactionApple record (if any)
	var existing models.TransactionApple
	findErr := mongoRepo.GetOne(
		"transactionApple",
		bson.M{"appleAppId": req.AppleAppId, "originalTransactionId": req.OriginalTransactionId},
		&existing,
	)

	// Generate App Store JWT
	jwtToken, err := utils.GenerateAppStoreJWT()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to generate JWT", err.Error())
		return
	}

	// Call Apple get-all-subscription-statuses
	data, err := services.FetchAllSubscriptionStatuses(jwtToken, req.OriginalTransactionId)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Failed to fetch subscription statuses", err.Error())
		return
	}

	var statusResp models.StatusResponse
	if err := json.Unmarshal(data, &statusResp); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to parse subscription statuses", err.Error())
		return
	}

	// Extract a representative status (pick first group -> first lastTransaction)
	var statusCode int32
	if len(statusResp.Data) > 0 && len(statusResp.Data[0].LastTransactions) > 0 {
		statusCode = statusResp.Data[0].LastTransactions[0].Status
	}

	// Map status to text
	statusText := mapStatusText(statusCode)

	// Decide active (not expired) – treat Active(1) and Grace Period(4) as active
	active := statusCode == 1 || statusCode == 4

	now := time.Now()

	// Prepare upsert doc
	doc := models.TransactionApple{
		AppleAppId:            req.AppleAppId,
		OriginalTransactionId: req.OriginalTransactionId,
		Status:                statusCode,
		StatusText:            statusText,
		BundleId:              statusResp.BundleId,
		Environment:           statusResp.Environment,
		UpdatedAt:             now,
	}

	if findErr == nil {
		// Update existing
		doc.ID = existing.ID
		if err := mongoRepo.Update(
			"transactionApple",
			bson.M{"_id": existing.ID},
			&doc,
		); err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to update transactionApple", err.Error())
			return
		}
	} else if !errors.Is(findErr, mongo.ErrNoDocuments) {
		respondWithError(ctx, http.StatusInternalServerError, "Database error while checking existing record", findErr.Error())
		return
	} else {
		// Create new
		if err := mongoRepo.Save("transactionApple", &doc); err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "Failed to save transactionApple", err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"exists":     findErr == nil,
		"status":     statusCode,
		"statusText": statusText,
		"active":     active,
		"bundleId":   statusResp.BundleId,
		"env":        statusResp.Environment,
	})
}

func mapStatusText(s int32) string {
	switch s {
	case 1:
		return "Active"
	case 2:
		return "Expired"
	case 3:
		return "Billing Retry"
	case 4:
		return "Billing Grace Period"
	case 5:
		return "Revoked"
	default:
		return "Unknown"
	}
}
