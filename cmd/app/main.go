package main

import (
	"database/sql"
	"log"

	"github.com/NewLeonardooliv/gateway-payment/internal/config"
	"github.com/NewLeonardooliv/gateway-payment/internal/repository"
	"github.com/NewLeonardooliv/gateway-payment/internal/service"
	"github.com/NewLeonardooliv/gateway-payment/internal/shared"
	"github.com/NewLeonardooliv/gateway-payment/internal/web/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := config.GetConnectionDatabase()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, *accountService)

	port := shared.GetEnv("HTTP_PORT", "8080")

	server := server.NewServer(accountService, invoiceService, port)
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server", err)
	}
}
