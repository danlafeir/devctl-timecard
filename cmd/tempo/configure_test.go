package tempo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.json")
	cfg := Config{TempoAPIToken: "sometoken123"}

	// Save config
	f, err := os.Create(configFile)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		t.Fatalf("Failed to encode config: %v", err)
	}

	// Load config
	f2, err := os.Open(configFile)
	if err != nil {
		t.Fatalf("Failed to open config file: %v", err)
	}
	defer f2.Close()
	var loaded Config
	if err := json.NewDecoder(f2).Decode(&loaded); err != nil {
		t.Fatalf("Failed to decode config: %v", err)
	}

	if loaded.TempoAPIToken != cfg.TempoAPIToken {
		t.Errorf("Expected token %q, got %q", cfg.TempoAPIToken, loaded.TempoAPIToken)
	}
}
