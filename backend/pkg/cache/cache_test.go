package cache

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSaveAndGetCache(t *testing.T) {
	tempFile := "test_cache.json"
	os.Setenv("CACHE_FILE_PATH", tempFile)
	defer os.Remove(tempFile)

	testData := map[string]string{"foo": "bar"}
	err := SaveCache("/test", testData, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	cached := GetCachedRoute("/test")
	if cached == nil {
		t.Fatal("Expected cached data, got nil")
	}

	var result map[string]string
	err = json.Unmarshal(cached, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal cached data: %v", err)
	}

	if result["foo"] != "bar" {
		t.Errorf("Expected bar, got %s", result["foo"])
	}
}

func TestCacheExpiry(t *testing.T) {
	tempFile := "test_expiry_cache.json"
	os.Setenv("CACHE_FILE_PATH", tempFile)
	defer os.Remove(tempFile)

	testData := map[string]string{"foo": "bar"}
	err := SaveCache("/expired", testData, -1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	cached := GetCachedRoute("/expired")
	if cached != nil {
		t.Fatal("Expected nil for expired cache, got data")
	}
}

func TestMiddleware(t *testing.T) {
	tempFile := "test_middleware_cache.json"
	os.Setenv("CACHE_FILE_PATH", tempFile)
	defer os.Remove(tempFile)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware(1 * time.Minute))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// First request - Cache Miss
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Errorf("Expected 200, got %d", w1.Code)
	}
	if w1.Header().Get("X-Cache") == "HIT" {
		t.Error("Expected cache MISS on first request")
	}

	// Second request - Cache Hit
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("Expected 200, got %d", w2.Code)
	}
	if w2.Header().Get("X-Cache") != "HIT" {
		t.Error("Expected cache HIT on second request")
	}
}

func TestGetCachedRouteErrors(t *testing.T) {
	os.Setenv("CACHE_FILE_PATH", "non_existent.json")
	if GetCachedRoute("/any") != nil {
		t.Error("Expected nil for non-existent file")
	}

	tempFile := "corrupt_cache.json"
	os.Setenv("CACHE_FILE_PATH", tempFile)
	os.WriteFile(tempFile, []byte("invalid json"), 0644)
	defer os.Remove(tempFile)

	if GetCachedRoute("/any") != nil {
		t.Error("Expected nil for corrupt json")
	}
}
