package db

import (
	"DeliFood/backend/pkg/logger"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDB(cfg Config, log *logger.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port,
		cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	log.Info("Attempting to connect to database", map[string]interface{}{
		"host":   cfg.Host,
		"port":   cfg.Port,
		"dbname": cfg.DBName,
	})

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Failed to open database", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		log.Error("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Info("Successfully connected to database", nil)
	return db, nil
}

func LoadConfigFromEnv(log *logger.Logger) Config {
	log.Info("Loading database configuration from environment variables", nil)
	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     getenvAsInt("DB_PORT", 5432),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

func getenvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
