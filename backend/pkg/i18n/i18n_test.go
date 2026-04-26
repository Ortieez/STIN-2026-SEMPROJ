package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetLanguage(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"en", "en"},
		{"cs", "cs"},
		{"cs-CZ", "cs"},
		{"fr", "en"}, // Fallback
		{"", "en"},   // Fallback
	}

	for _, tt := range tests {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Accept-Language", tt.header)

		lang := GetLanguage(c)
		if lang != tt.expected {
			t.Errorf("For header %s, expected %s, got %s", tt.header, tt.expected, lang)
		}
	}
}

func TestT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		lang     string
		key      string
		expected string
	}{
		{"en", "unauthorized", "Unauthorized"},
		{"cs", "unauthorized", "Neautorizovaný přístup"},
		{"en", "settings_saved", "Settings saved"},
		{"cs", "settings_saved", "Nastavení uloženo"},
		{"en", "non_existent_key", "non_existent_key"},
	}

	for _, tt := range tests {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Accept-Language", tt.lang)

		res := T(c, tt.key)
		if res != tt.expected {
			t.Errorf("For lang %s and key %s, expected %s, got %s", tt.lang, tt.key, tt.expected, res)
		}
	}
}
