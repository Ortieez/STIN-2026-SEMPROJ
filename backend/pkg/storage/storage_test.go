package storage

import (
	"os"
	"testing"
)

func TestStorageSettings(t *testing.T) {
	tempSettings := "test_settings_unit.json"
	os.Setenv("USER_SETTINGS_PATH", tempSettings)
	defer os.Remove(tempSettings)

	s := NewStorage()

	// Initial settings (defaults)
	settings, err := s.GetSettings()
	if err != nil {
		t.Fatalf("Error getting settings: %v", err)
	}
	if settings.BaseCurrency != "EUR" {
		t.Errorf("Expected default EUR, got %s", settings.BaseCurrency)
	}

	// Save and Load
	newSettings := UserSettings{
		BaseCurrency:       "CZK",
		SelectedCurrencies: []string{"USD", "EUR"},
	}
	err = s.SaveSettings(newSettings)
	if err != nil {
		t.Fatalf("Error saving settings: %v", err)
	}

	loaded, err := s.GetSettings()
	if err != nil {
		t.Fatalf("Error loading settings: %v", err)
	}
	if loaded.BaseCurrency != "CZK" || len(loaded.SelectedCurrencies) != 2 {
		t.Errorf("Settings mismatch: %+v", loaded)
	}
}

func TestStorageLogs(t *testing.T) {
	tempLog := "test_app_unit.log"
	os.Setenv("LOG_FILE_PATH", tempLog)
	defer os.Remove(tempLog)

	s := NewStorage()
	s.Log("INFO", "Test message")

	if _, err := os.Stat(tempLog); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}
