package routers

import (
	"simvizlab-backend/controllers/user"

	"github.com/gin-gonic/gin"
)

// UserRoutes registers all user routes with JWT authentication
func UserRoutes(rg *gin.RouterGroup) {
	// Remove the additional /user group since it's already grouped in index.go
	rg.GET("/", user.GetAllUsers)
	// rg.GET("/:id", user.GetUserByID)
	rg.POST("/", user.CreateUser)
	rg.PUT("/:id", user.UpdateUser)

	rg.GET("/transaction",user.GetTransaction)
	// rg.DELETE("/:id", user.DeleteUser)
}
