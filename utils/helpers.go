package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, err error, message string) {
	errorResponse := ErrorResponse{
		Status:  statusCode,
		Message: message,
	}

	log.Printf("Error: %v", err)

	c.JSON(statusCode, errorResponse)
}
