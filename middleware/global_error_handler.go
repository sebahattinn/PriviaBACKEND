package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic occurred: %v", r)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "Something went wrong. Please try again later.",
				})
			}
		}()

		c.Next()
	}
}
