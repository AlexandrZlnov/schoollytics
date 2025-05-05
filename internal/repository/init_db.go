package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	//"log"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// содзает подключение к базе данных
// возвращает *sql.DB или error
func InitDB() (*sql.DB, error) {
	// загружаем конфигурацию из .env
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	// connStr := "host=localhost port=5432 user=postgres password=Zel408 dbname=test4 sslmode=disable"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	fmt.Println("connStr in initDB ----> ", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, fmt.Errorf("ошибка подключения к PostgreSql: %v", err)
	}

	if err := db.Ping(); err != nil {
		return db, fmt.Errorf("ошибка Ping DB: %v", err)
	}

	return db, nil
}

func loadConfig() (*domain.Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
	}
	return &domain.Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}, nil
}

// вспомогательная функция к func loadConfig()
// присваивает значение из .env или дефолнтное
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
