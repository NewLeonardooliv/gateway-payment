package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/NewLeonardooliv/gateway-payment/internal/repository"
	"github.com/NewLeonardooliv/gateway-payment/internal/service"
	"github.com/NewLeonardooliv/gateway-payment/internal/web/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%s port=%s port=%s port=%s port=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	port := getEnv("HTTP_PORT", "8080")

	server := server.NewServer(accountService, port)
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server", err)
	}
}
