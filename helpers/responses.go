package helpers

import "github.com/gin-gonic/gin"

type HTTPResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func GetHTTPResponse(statusCode int, context *gin.Context, response HTTPResponse) {
	context.JSON(statusCode, gin.H{
		"result":  response.Result,
		"message": response.Message,
		"data":    response.Data,
	})
}
