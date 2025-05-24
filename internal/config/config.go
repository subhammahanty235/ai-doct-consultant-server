package config

import (
	"os"
)

type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	JWTSecret    string
	GeminiAPIKey string
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	S3Bucket     string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "ai_doctor_db"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY", ""),
		AWSSecretKey: getEnv("AWS_SECRET_KEY", ""),
		S3Bucket:     getEnv("S3_BUCKET", "ai-doctor-images"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
