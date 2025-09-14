package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"simvizlab-backend/models"
	mongoRepo "simvizlab-backend/repository/mongo"
	"simvizlab-backend/services"
	"simvizlab-backend/utils"

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
