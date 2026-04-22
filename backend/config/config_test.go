package config

import (
	"testing"
)

// test func LoadConfig()
// run test by "go test -v"

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()

	// test expect to is empty string
	if cfg.GeminiAPIKey == "" {
		t.Error("Expected LineChannelAccessToken to not be empty, but got empty string")
	}

	if cfg.GeminiAPIKey == "" {
		t.Error("Expected GeminiAPIKey to not be empty, but got empty string")
	}
}
