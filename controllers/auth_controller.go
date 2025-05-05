package controllers

import (
	"net/http"
	"priviatodolist/middleware"
	"priviatodolist/mockdb"

	"github.com/gin-gonic/gin"
)

// Login handles user authentication and JWT token generation
// @Summary User login
// @Description Authenticate with username and password to receive a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Login Credentials"
// @Success 200 "Authentication token"
// @Failure 400 "Invalid JSON or credentials"
// @Failure 401 "User role or ID not found"
// @Failure 500 "Token generation failed"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	role := mockdb.UserRoles[loginData.Username]
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	userID := mockdb.UserIDs[loginData.Username]
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	token, err := middleware.GenerateToken(userID, loginData.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
