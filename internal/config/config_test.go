package config

import "testing"

func TestLoad_missingURL(t *testing.T) {
	// Arrange
	t.Setenv("MINIFLUX_URL", "")
	t.Setenv("MINIFLUX_API_KEY", "secret")

	// Act
	_, err := Load()

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_missingAPIKey(t *testing.T) {
	// Arrange
	t.Setenv("MINIFLUX_URL", "https://miniflux.example")
	t.Setenv("MINIFLUX_API_KEY", "")

	// Act
	_, err := Load()

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_ok(t *testing.T) {
	// Arrange
	t.Setenv("MINIFLUX_URL", "https://miniflux.example")
	t.Setenv("MINIFLUX_API_KEY", "secret")

	// Act
	cfg, err := Load()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://miniflux.example" || cfg.APIKey != "secret" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
