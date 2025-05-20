package config

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")
	val := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", val)
}

func TestGetEnv_Fallback(t *testing.T) {
	os.Unsetenv("TEST_KEY_NOT_SET")
	val := getEnv("TEST_KEY_NOT_SET", "default")
	assert.Equal(t, "default", val)
}

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DB_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("AUTH_SERVICE_URL")

	cfg := Load()
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/forum?sslmode=disable", cfg.DBURL)
	assert.Equal(t, "default-secret-key", cfg.JWTSecret)
	assert.Equal(t, "http://localhost:8081", cfg.AuthServiceURL)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "1234")
	os.Setenv("DB_URL", "dburl")
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("AUTH_SERVICE_URL", "http://auth")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DB_URL")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("AUTH_SERVICE_URL")

	cfg := Load()
	assert.Equal(t, 1234, cfg.Port)
	assert.Equal(t, "dburl", cfg.DBURL)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, "http://auth", cfg.AuthServiceURL)
} 