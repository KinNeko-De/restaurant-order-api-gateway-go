package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type orderCreatedResponse struct {
	Id uuid.UUID `json:"id"`
}

type orderResponse struct {
	Id      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}

func CreateOrder(context *gin.Context) {
	var response orderCreatedResponse

	response.Id = uuid.New()
	getOrderUri := GetHostSchema(context) + context.Request.Host + "/order/" + response.Id.String()
	context.Writer.Header().Set("Location", getOrderUri)
	context.JSON(http.StatusCreated, response)
}

const ParamOrderId string = "orderId"

func GetOrder(context *gin.Context) {
	paramId := context.Param(ParamOrderId)
	requestId, err := uuid.Parse(paramId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response orderResponse
	response.Id = requestId
	response.Version = 1

	context.JSON(http.StatusOK, response)
}

func GetHostSchema(context *gin.Context) string {
	scheme := "http"
	if context.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://"
}
