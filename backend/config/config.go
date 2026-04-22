package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// struct format from .env file and .env.example
type Config struct {
	LineChannelAccessToken string
	GeminiAPIKey           string
}

// Function for load env
func LoadConfig() Config {

	// godotenv.Load("path of .ennv")
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error: No .env file found. Falling back to system environment variables.")
	}

	conf := Config{
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
	}
	return conf
}
