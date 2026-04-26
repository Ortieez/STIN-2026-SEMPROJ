package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func Middleware() gin.HandlerFunc {
	_ = godotenv.Load()
	expectedToken := os.Getenv("AUTH_TOKEN")
	if expectedToken == "" {
		// Default token for dev is hash of admin:password123
		expectedToken = hashString("admin:password123")
	}

	return func(c *gin.Context) {
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
	Username string `json:"username"` // Expected to be SHA256 hash from client
	Password string `json:"password"` // Expected to be SHA256 hash from client
}

func LoginHandler(c *gin.Context) {
	_ = godotenv.Load()
	
	// Plain text credentials from .env
	plainUsername := os.Getenv("LOGIN_USERNAME")
	plainPassword := os.Getenv("LOGIN_PASSWORD")
	token := os.Getenv("AUTH_TOKEN")

	// Defaults for development
	if plainUsername == "" {
		plainUsername = "admin"
	}
	if plainPassword == "" {
		plainPassword = "password123"
	}
	if token == "" {
		token = hashString("admin:password123")
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Server-side hashing of plain credentials from .env
	expectedUserHash := hashString(plainUsername)
	expectedPassHash := hashString(plainPassword)

	// Compare incoming hashes with locally generated hashes
	if req.Username == expectedUserHash && req.Password == expectedPassHash {
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
