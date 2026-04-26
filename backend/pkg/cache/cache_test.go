package cache

import (
	"encoding/json"
	"os"
	"testing"
	"time"
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
	// Save with negative duration (already expired)
	err := SaveCache("/expired", testData, -1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	cached := GetCachedRoute("/expired")
	if cached != nil {
		t.Fatal("Expected nil for expired cache, got data")
	}
}
