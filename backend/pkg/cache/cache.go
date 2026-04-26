package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type CacheEntry struct {
	ExpiryDate time.Time       `json:"expiryDate"`
	Data       json.RawMessage `json:"data"`
}

type FullCache map[string]CacheEntry

func getCacheFilename() string {
	cacheFileName := os.Getenv("CACHE_FILE_PATH")
	if cacheFileName != "" {
		return cacheFileName
	}

	_ = godotenv.Load()
	return os.Getenv("CACHE_FILE_PATH")
}

func GetCachedRoute(route string) []byte {
	filename := getCacheFilename()

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	var allCache FullCache
	if err := json.Unmarshal(fileData, &allCache); err != nil {
		return nil
	}

	entry, exists := allCache[route]
	if !exists {
		return nil
	}

	if time.Now().After(entry.ExpiryDate) {
		return nil
	}

	return entry.Data
}

func SaveCache(route string, data interface{}, duration time.Duration) error {
	filename := getCacheFilename()
	allCache := make(FullCache)

	fileData, err := os.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(fileData, &allCache)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	allCache[route] = CacheEntry{
		ExpiryDate: time.Now().Add(duration),
		Data:       jsonData,
	}

	fileBytes, err := json.MarshalIndent(allCache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, fileBytes, 0644)
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Middleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.RequestURI

		cachedData := GetCachedRoute(key)
		if cachedData != nil {
			c.Header("X-Cache", "HIT")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(200)
			c.Writer.Write(cachedData)
			return
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		if c.Writer.Status() == 200 && w.body.Len() > 0 {
			err := SaveCache(key, json.RawMessage(w.body.Bytes()), duration)
			if err != nil {
				fmt.Printf("Cache Save Error: %v\n", err)
			}
		}
	}
}
