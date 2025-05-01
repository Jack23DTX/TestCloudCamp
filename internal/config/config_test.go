package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный конфиг-файл
	configContent := `
port: ":8080"
backends:
  - "http://localhost:8081"
rate_limit:
  rps: 5
  burst: 10
`
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Тест загрузки корректного конфига
	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}
	if cfg.Port != ":8080" {
		t.Errorf("Expected port :8080, got %s", cfg.Port)
	}
}

func TestLoadConfig_FileNotExists(t *testing.T) {
	_, err := LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
