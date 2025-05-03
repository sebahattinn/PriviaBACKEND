package controllers

import (
	"net/http"
	"priviatodolist/middleware"
	"priviatodolist/mockdb"

	"github.com/gin-gonic/gin"
)

// Login kullanıcı login işlemi ve JWT oluşturma
// @Summary Kullanıcı girişi yapar
// @Description Kullanıcı adı ve şifre ile giriş yaparak JWT token alır
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Kullanıcı Girişi"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "Hatalı JSON"
// @Failure 401 {object} map[string]string "Geçersiz kimlik bilgisi"
// @Failure 500 {object} map[string]string "Token üretimi başarısız"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	pass, ok := mockdb.Users[loginData.Username]
	if !ok || pass != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	role := mockdb.UserRoles[loginData.Username]

	token, err := middleware.GenerateToken(loginData.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
