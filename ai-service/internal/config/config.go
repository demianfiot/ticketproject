package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		GRPCPort string `mapstructure:"grpc_port"`
	} `mapstructure:"server"`

	OpenAI struct {
		Model          string `mapstructure:"model"`
		BaseURL        string `mapstructure:"base_url"`
		TimeoutSeconds int    `mapstructure:"timeout_seconds"`
	} `mapstructure:"openai"`

	OpenAIAPIKey string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	_ = godotenv.Load(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	return &cfg, nil
}
