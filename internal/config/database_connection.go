package config

import (
	"fmt"

	"github.com/NewLeonardooliv/gateway-payment/internal/shared"
)

func GetConnectionDatabase() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		shared.GetEnv("DB_HOST", "db"),
		shared.GetEnv("DB_PORT", "5432"),
		shared.GetEnv("DB_USER", "postgres"),
		shared.GetEnv("DB_PASSWORD", "postgres"),
		shared.GetEnv("DB_NAME", "gateway"),
		shared.GetEnv("DB_SSL_MODE", "disable"),
	)
}

func GetConnectionStringDatabase() string {
	user := shared.GetEnv("DB_USER", "")
	password := shared.GetEnv("DB_PASSWORD", "")
	host := shared.GetEnv("DB_HOST", "")
	port := shared.GetEnv("DB_PORT", "")
	dbname := shared.GetEnv("DB_NAME", "")
	sslmode := shared.GetEnv("DB_SSL_MODE", "")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	return connStr
}
