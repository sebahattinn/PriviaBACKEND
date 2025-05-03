package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY is not set in the environment")
	}

	jwtSecret = []byte(secret)
	log.Println("Secret key loaded successfully") // Geliştiriciye özel log
}

// JWT Middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("No token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// "Bearer " varsa temizle
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		log.Println("Parsing token:", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Invalid signing method")
				return nil, http.ErrNotSupported
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		log.Println("Token claims parsed successfully")

		// 🔧 Burada artık userID int olarak alınıyor, arayüzde string olabilir o yüzden dönüştür
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			log.Println("userID not found or not a number")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in token"})
			c.Abort()
			return
		}

		// Token'den "userID" alınıp "ownerID" olarak ayarlanır
		c.Set("ownerID", int(userIDFloat))
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// Token Oluşturucu
func GenerateToken(userID int, username, role string) (string, error) {
	log.Println("Generating token for userID:", userID, "username:", username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	log.Println("Token generated successfully")
	return signedToken, nil
}
