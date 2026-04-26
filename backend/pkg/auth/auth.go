package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Middleware() gin.HandlerFunc {
	_ = godotenv.Load()
	expectedToken := os.Getenv("AUTH_TOKEN")
	if expectedToken == "" {
		expectedToken = "secret-token"
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
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	_ = godotenv.Load()
	expectedUsername := os.Getenv("LOGIN_USERNAME")
	expectedPassword := os.Getenv("LOGIN_PASSWORD")
	token := os.Getenv("AUTH_TOKEN")

	if expectedUsername == "" {
		expectedUsername = "admin"
	}
	if expectedPassword == "" {
		expectedPassword = "password123"
	}
	if token == "" {
		token = "secret-token"
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username == expectedUsername && req.Password == expectedPassword {
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
