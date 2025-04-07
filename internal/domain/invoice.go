package domain

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID             string
	AccountID      string
	Amount         float64
	Status         string
	Description    string
	PaymentType    string
	CardLastDigits string
	mu             sync.RWMutex
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

func NewInvoice(account *Account, amount float64, description string, paymentType string, cardLastDigits string) *Invoice {
	invoice := &Invoice{
		ID:             uuid.New().String(),
		AccountID:      account.ID,
		Amount:         amount,
		Status:         "pending",
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: cardLastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return invoice
}

func (invoice *Invoice) ChangeStatus(status string) {
	invoice.mu.Lock()
	defer invoice.mu.Unlock()

	invoice.Status = status
	invoice.UpdatedAt = time.Now()
}
