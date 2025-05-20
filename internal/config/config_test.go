package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "полная конфигурация",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				DBName:   "forum",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=forum sslmode=disable",
		},
		{
			name: "минимальная конфигурация",
			config: DatabaseConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				User:     "user",
				Password: "",
				DBName:   "test",
				SSLMode:  "require",
			},
			expected: "host=127.0.0.1 port=5432 user=user password= dbname=test sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.GetDSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл конфигурации
	configContent := `
database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: forum
  sslmode: disable
  driver: postgres
  max_open_conns: 25
  max_idle_conns: 25
  conn_max_lifetime: 5m
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	// Тестируем успешную загрузку конфигурации
	config, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Проверяем значения
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "password", config.Database.Password)
	assert.Equal(t, "forum", config.Database.DBName)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "postgres", config.Database.Driver)
	assert.Equal(t, 25, config.Database.MaxOpenConns)
	assert.Equal(t, 25, config.Database.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, config.Database.ConnMaxLifetime)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	config, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "ошибка чтения файла конфигурации")
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Создаем временный файл с некорректным YAML
	tmpFile, err := os.CreateTemp("", "invalid-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("invalid: yaml: content: [")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	config, err := LoadConfig(tmpFile.Name())
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "ошибка парсинга конфигурации")
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	// Создаем пустой временный файл
	tmpFile, err := os.CreateTemp("", "empty-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	require.NoError(t, tmpFile.Close())

	// Записываем пустой YAML документ
	err = os.WriteFile(tmpFile.Name(), []byte("{}"), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.NotNil(t, config.Database)
} 