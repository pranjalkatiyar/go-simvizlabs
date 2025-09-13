package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simvizlab-backend/models"
	"simvizlab-backend/repository/mongo"
	"simvizlab-backend/services"
	"simvizlab-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(ctx *gin.Context) {
	var users []*models.User
	err := mongo.Get("users", &users)

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
	err := mongo.GetOne("users", nil, &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context) {

	// 1. Bind request
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AppleAppId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "AppleAppId must be a positive integer"})
		return
	}

	//check if user already exist or not ?
	var existingUser models.User
	err := mongo.GetOne("users", bson.M{"appleAppId": req.AppleAppId}, &existingUser)

	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User with this AppleAppId already exists", "existUser": true})
		return
	}

	jwtToken, err := utils.GenerateAppStoreJWT()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	// 2. Call Apple Transaction API
	data, err := services.FetchTransaction(jwtToken, req.TransactionId)
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

	fmt.Println("results:", results)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction data"})
		return
	}

	user := models.User{
		AppleAppId: req.AppleAppId,
		
	}

	// err := mongo.Save("users", &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := mongo.Update("users", map[string]interface{}{"_id": user.ID}, &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
