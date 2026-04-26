package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginHandler(t *testing.T) {
	// Plain credentials in .env
	os.Setenv("LOGIN_USERNAME", "admin")
	os.Setenv("LOGIN_PASSWORD", "password123")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", LoginHandler)

	// Valid login: client sends hashes
	w := httptest.NewRecorder()
	hashedUser := hashString("admin")
	hashedPass := hashString("password123")
	
	body, _ := json.Marshal(LoginRequest{Username: hashedUser, Password: hashedPass})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Invalid login: wrong hashes
	w = httptest.NewRecorder()
	body, _ = json.Marshal(LoginRequest{Username: "wronghash", Password: "wronghash"})
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestMiddleware(t *testing.T) {
	os.Setenv("LOGIN_USERNAME", "admin")
	os.Setenv("LOGIN_PASSWORD", "password123")
	expectedToken := hashString(fmt.Sprintf("%s:%s", "admin", "password123"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	router.GET("/protected", func(c *gin.Context) { c.Status(200) })

	// Valid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", expectedToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Invalid token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "invalid")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}
