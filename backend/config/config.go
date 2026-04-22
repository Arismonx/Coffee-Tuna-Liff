package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LineChannelAccessToken string
	GeminiAPIKey           string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file!")
	}
	conf := Config{
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
	}
	return conf
}
