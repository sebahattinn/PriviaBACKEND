package middleware

import (
	"net/http"

	"priviatodolist/mockdb"

	"github.com/gin-gonic/gin"
)

// AdminOnly middleware'ı sadece admin kullanıcılarına erişim izni verir.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		role, exists := mockdb.UserRoles[userID.(string)]

		// Eğer kullanıcı admin değilse, erişim reddedilsin
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin only route"})
			c.Abort()
			return
		}

		c.Next()
	}
}
