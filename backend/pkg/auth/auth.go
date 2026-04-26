package auth

import (
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

// getAuthCredentials retrieves values from .env with minimal hardcoded fallbacks
func getAuthCredentials() (string, string, string) {
	_ = godotenv.Load()
	user := os.Getenv("LOGIN_USERNAME")
	pass := os.Getenv("LOGIN_PASSWORD")
	token := os.Getenv("AUTH_TOKEN")

	// If token is missing, generate it dynamically from current user/pass
	if token == "" {
		u := user
		if u == "" {
			u = "admin"
		}
		p := pass
		if p == "" {
			p = "password123"
		}
		token = hashString(fmt.Sprintf("%s:%s", u, p))
	}

	return user, pass, token
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, expectedToken := getAuthCredentials()
		
		token := c.GetHeader("Authorization")
		if token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

	// Use defaults if .env is completely empty
	if envUser == "" {
		envUser = "admin"
	}
	if envPass == "" {
		envPass = "password123"
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	expectedUserHash := hashString(envUser)
	expectedPassHash := hashString(envPass)

	if req.Username == expectedUserHash && req.Password == expectedPassHash {
		c.JSON(http.StatusOK, gin.H{"token": expectedToken})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
