package bootstrap

import (
	"log"

	"database/sql"

	"github.com/NewLeonardooliv/gateway-payment/internal/config"
	account_repository "github.com/NewLeonardooliv/gateway-payment/internal/repository/implementations/account"
	invoice_repository "github.com/NewLeonardooliv/gateway-payment/internal/repository/implementations/invoice"
	"github.com/NewLeonardooliv/gateway-payment/internal/service"
	"github.com/NewLeonardooliv/gateway-payment/internal/shared"
	"github.com/NewLeonardooliv/gateway-payment/internal/web/server"
)

func LoadWeb() {
	connectionString := config.GetConnectionDatabase()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	defer db.Close()

	accountRepository := account_repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	interInvoiceRepository := invoice_repository.NewInterInvoiceRepository(
		shared.GetEnv("INTERBANK_CLIENT_ID", ""),
		shared.GetEnv("INTERBANK_CLIENT_SECRET", ""),
	)

	invoiceService := service.NewInvoiceService(interInvoiceRepository, *accountService)

	port := shared.GetEnv("HTTP_PORT", "8080")

	server := server.NewServer(accountService, invoiceService, port)
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server", err)
	}
}
