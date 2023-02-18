package main

import (
	"github.com/KinNeko-De/restaurant-order-api-gateway-go/order"
	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()

	_ = router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/order", order.CreateOrder)
	router.GET("/order/:"+order.GetOrderParamOrderId, order.GetOrder)
	return router
}
