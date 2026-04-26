package auth

import (
	"backend/pkg/i18n"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func getAuthCredentials() (string, string, string) {
	_ = godotenv.Load()
	user := os.Getenv("LOGIN_USERNAME")
	pass := os.Getenv("LOGIN_PASSWORD")
	token := ""

	u := user
	p := pass
	token = hashString(fmt.Sprintf("%s:%s", u, p))

	return user, pass, token
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, expectedToken := getAuthCredentials()

		token := c.GetHeader("Authorization")
		if token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(c, "unauthorized")})
			c.Abort()
			return
		}
		c.Next()
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	envUser, envPass, expectedToken := getAuthCredentials()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(c, "invalid_request")})
		return
	}

	expectedUserHash := hashString(envUser)
	expectedPassHash := hashString(envPass)

	if req.Username == expectedUserHash && req.Password == expectedPassHash {
		c.JSON(http.StatusOK, gin.H{"token": expectedToken})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.T(c, "invalid_credentials")})
	}
}
