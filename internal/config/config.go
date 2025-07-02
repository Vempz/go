package config

import (
    "log"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Port            string
    Environment     string
    RedisAddr       string
    RedisPassword   string
    SessionSecret   string
    SessionMaxAge   time.Duration
    JWTSecret       string
    JWTExpiresIn    time.Duration
    AllowedOrigins  []string
    BCryptCost      int
}

func Load() *Config {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    port := getEnv("PORT", "3000")
    environment := getEnv("NODE_ENV", "development")
    redisAddr := getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379")
    redisPassword := getEnv("REDIS_PASSWORD", "")
    sessionSecret := getEnv("SESSION_SECRET", "your-secret-key")
    
    sessionMaxAge, _ := strconv.Atoi(getEnv("SESSION_MAX_AGE", "86400"))
    jwtExpiresHours, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))
    bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "12"))

    allowedOrigins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ",")

    return &Config{
        Port:            port,
        Environment:     environment,
        RedisAddr:       redisAddr,
        RedisPassword:   redisPassword,
        SessionSecret:   sessionSecret,
        SessionMaxAge:   time.Duration(sessionMaxAge) * time.Second,
        JWTSecret:       getEnv("JWT_SECRET", "jwt-secret"),
        JWTExpiresIn:    time.Duration(jwtExpiresHours) * time.Hour,
        AllowedOrigins:  allowedOrigins,
        BCryptCost:      bcryptCost,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
