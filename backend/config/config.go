package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LineChannelAccessToken string
	GeminiAPIKey           string
}

// Function for load env
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file!")
	}

	conf := Config{
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
	}
	return conf
}
