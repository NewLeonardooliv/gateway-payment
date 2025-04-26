package main

import (
	"log"

	"github.com/NewLeonardooliv/gateway-payment/database"
	"github.com/NewLeonardooliv/gateway-payment/internal/bootstrap"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Migrate()
	bootstrap.LoadWeb()
}
