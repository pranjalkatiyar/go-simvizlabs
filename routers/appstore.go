package routers

import (
	"github.com/gin-gonic/gin"
)

func AppStoreRoutes(route *gin.RouterGroup) {
	route.POST("/transaction", controller.GetTransactionInfo)
}
