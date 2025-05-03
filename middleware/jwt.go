package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header'ını alıyoruz
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("No token provided") // Log: Token sağlanmamış
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Token'ı parse ediyoruz
		log.Println("Parsing token:", tokenString) // Log: Token çözülmeye çalışılıyor
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Invalid signing method") // Log: Geçersiz imza metodu
				return nil, http.ErrNotSupported
			}
			return []byte("your-secret-key"), nil // Burada secret key'inizi kullandığınızı unutmayın
		})

		// Token geçersizse hata ver
		if err != nil || !token.Valid {
			log.Println("Invalid token:", err) // Log: Geçersiz token hatası
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Claims'leri alıyoruz
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid token claims") // Log: Geçersiz token claim'leri
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Claims'ten username ve role alınıyor
		log.Println("Token claims parsed successfully") // Log: Claim'ler başarıyla çözüldü
		c.Set("userID", claims["username"])             // Gerekirse burada claims["userID"] olarak değiştir
		c.Set("role", claims["role"])

		c.Next()
	}
}

func GenerateToken(username, role string) (string, error) {
	// Token oluşturuluyor
	log.Println("Generating token for username:", username) // Log: Token oluşturuluyor
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// Token signed ve string olarak döndürülüyor
	signedToken, err := token.SignedString([]byte("PriviadaStajYapmakIstiyorumLutfenBeniStajaAlinHemGuzelKodYaziyorum"))
	if err != nil {
		log.Println("Error signing token:", err) // Log: Token imzalanırken hata oluştu
		return "", err
	}

	log.Println("Token generated successfully") // Log: Token başarıyla oluşturuldu
	return signedToken, nil
}
