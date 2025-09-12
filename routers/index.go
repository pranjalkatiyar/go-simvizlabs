package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	// Create an api group for all routes
	api := route.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
		})

		// Add all other routes within the api group
		UserRoutes(api.Group("/user"))
		AppStoreRoutes(api.Group("/appstore"))
	}
}
