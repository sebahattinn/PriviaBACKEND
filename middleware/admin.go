package middleware

import (
	"github.com/gin-gonic/gin"
)

// AdminOnly middleware'ı sadece admin kullanıcılarına erişim izni verir.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(403, gin.H{"error": "Yalnızca admin erişebilir"})
			c.Abort()
			return
		}

		c.Next()
	}
}
