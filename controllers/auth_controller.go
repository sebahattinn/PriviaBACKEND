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

	// JSON verilerini kontrol etme
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Kullanıcı adı ve şifreyi kontrol etme
	pass, ok := mockdb.Users[loginData.Username]
	if !ok || pass != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Kullanıcı rolünü al
	role := mockdb.UserRoles[loginData.Username]
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Kullanıcı ID'sini al
	userID := mockdb.UserIDs[loginData.Username]
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Token üretme
	token, err := middleware.GenerateToken(userID, loginData.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Başarılı giriş ve token dönüşü
	c.JSON(http.StatusOK, gin.H{"token": token})
}
