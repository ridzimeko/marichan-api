package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName         string
	AppEnv          string
	AppPort         string
	BaseURL         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	DBSSLMode       string
	APIKey          string
	UploadDir       string
	TempDir         string
	MaxUploadSizeMB int64

	PixivPHPSESSID string
	PixivUserAgent string
	PixivBaseURL   string

	GeminiAPIKey           string
	GeminiSystemPromptFile string

	GroqAPIKey             string
}

func LoadEnv() *Env {
	_ = godotenv.Load()

	env := &Env{
		AppName:                getEnv("APP_NAME", "marichan-api"),
		AppEnv:                 getEnv("APP_ENV", "development"),
		AppPort:                getEnv("APP_PORT", "8080"),
		BaseURL:                getEnv("BASE_URL", "http://localhost:8080"),
		DBHost:                 getEnv("DB_HOST", "127.0.0.1"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPass:                 getEnv("DB_PASS", "postgres"),
		DBName:                 getEnv("DB_NAME", "bot_api"),
		DBSSLMode:              getEnv("DB_SSLMODE", "disable"),
		APIKey:                 getEnv("API_KEY", "supersecretkey"),
		UploadDir:              getEnv("UPLOAD_DIR", "uploads"),
		TempDir:                getEnv("TEMP_DIR", "tmp"),
		MaxUploadSizeMB:        getEnvInt64("MAX_UPLOAD_SIZE_MB", 10),
		PixivPHPSESSID:         getEnv("PIXIV_PHPSESSID", ""),
		PixivUserAgent:         getEnv("PIXIV_USER_AGENT", "Mozilla/5.0"),
		PixivBaseURL:           getEnv("PIXIV_BASE_URL", "https://www.pixiv.net"),
		GeminiAPIKey:           getEnv("GEMINI_API_KEY", ""),
		GeminiSystemPromptFile: getEnv("GEMINI_SYSTEM_PROMPT_FILE", "prompts/chatbot_mari.txt"),
		GroqAPIKey:             getEnv("GROQ_API_KEY", ""),
	}

	log.Printf("loaded env: %s (%s)\n", env.AppName, env.AppEnv)
	return env
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getEnvInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return n
}
