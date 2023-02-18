package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type orderCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}

type orderResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}

func CreateOrder(context *gin.Context) {
	var response orderCreatedResponse

	response.ID = uuid.New()

	context.JSON(http.StatusCreated, response)
}

const GetOrderParamOrderId string = "orderId"

func GetOrder(context *gin.Context) {
	paramId := context.Param(GetOrderParamOrderId)
	requestId, err := uuid.Parse(paramId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response orderResponse
	response.ID = requestId
	response.Version = 1

	context.JSON(http.StatusOK, response)
}
