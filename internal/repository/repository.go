package repository

import "github.com/NewLeonardooliv/gateway-payment/internal/domain"

type AccountRepository interface {
	Save(account *domain.Account) error
	FindByAPIKey(apiKey string) (*domain.Account, error)
	FindByID(id string) (*domain.Account, error)
	UpdateBalance(account *domain.Account) error
}

type InvoiceRepository interface {
	Save(invoice *domain.Invoice) error
	FindByID(id string) (*domain.Invoice, error)
	FindByAccountID(accountID string) ([]*domain.Invoice, error)
	UpdateStatus(invoice *domain.Invoice) error
}
