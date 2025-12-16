package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode               string
	DBUrl              string
	AppPort            string
	ImageKitPrivateKey string
	ImageKitPublicKey  string
	ImageKitEndpoint   string
	ZaloMiniAppHost    string
	ZaloAppID          string
	ZaloAppPrivateKey  string
	ZaloAppSecret      string
}

var AppConfig *Config

func InitConfig() *Config {
	LoadEnv()
	AppConfig = &Config{
		Mode:               getEnv("APP_MODE", "development"),
		DBUrl:              getEnv("DB_URL", ""),
		AppPort:            getEnv("APP_PORT", "8080"),
		ImageKitPrivateKey: getEnv("IMAGEKIT_PRIVATE_KEY", ""),
		ImageKitPublicKey:  getEnv("IMAGEKIT_PUBLIC_KEY", ""),
		ImageKitEndpoint:   getEnv("IMAGEKIT_ENDPOINT_URL", ""),
		ZaloMiniAppHost:    getEnv("ZALO_MINIAPP_HOST", ""),
		ZaloAppPrivateKey:  getEnv("ZALO_MINIAPP_PRIVATE_KEY", ""),
		ZaloAppID:          getEnv("ZALO_APP_ID", ""),
		ZaloAppSecret:      getEnv("ZALO_APP_SECRET", ""),
	}

	return AppConfig
}

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	// Find .env file in the project root or current directory
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
		// Continue execution even if .env is not found
		// Variables might be set in the environment directly
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func IsProdEnv() bool {
	return AppConfig.Mode == "production"
}
