package routers

import (
	controller "simvizlab-backend/controllers/appstore"

	"github.com/gin-gonic/gin"
)

func AppStoreRoutes(route *gin.RouterGroup) {
	route.POST("/transaction", controller.GetTransactionInfo)
	route.POST("/history", controller.GetHistoryInfo)
}
