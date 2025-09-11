package user

import (
	"net/http"
	"simvizlab-backend/models"
	"simvizlab-backend/repository/mongo"

	"github.com/gin-gonic/gin"
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
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := mongo.Save("users", &user)

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
