package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly middleware'ı sadece admin kullanıcılarına erişim izni verir.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Burada admin kontrolü yapmamız lazım. Örnek olarak userID'yi kontrol ediyoruz.
		// Gerçek senaryoda veritabanından rolü kontrol etmelisin.
		if userID != 1 { // admin ID'sini sabit olarak 1 kabul ediyoruz
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin only route"})
			c.Abort()
			return
		}

		c.Next()
	}
}
