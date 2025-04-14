package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Invoice struct {
	ID             string
	AccountID      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type CreditCard struct {
	Number         string
	CVV            string
	ExpiryMonth    int
	ExpiryYear     int
	CardholderName string
}

func NewInvoice(accountID string, amount float64, description string, paymentType string, card CreditCard) (*Invoice, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	cardLastDigits := card.Number[len(card.Number)-4:]

	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Amount:         amount,
		Status:         "pending",
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: cardLastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (invoice *Invoice) Process(status Status) error {
	if invoice.Amount > 10000 {
		return nil
	}

	randomSource := rand.New(rand.NewSource(time.Now().Unix()))
	var newStatus Status

	newStatus = StatusRejected
	if randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	}

	invoice.UpdateStatus(newStatus)

	return nil
}

func (invoice *Invoice) UpdateStatus(newStatus Status) error {
	if invoice.Status != StatusPending {
		return ErrInvalidStatus
	}

	invoice.Status = newStatus
	invoice.UpdatedAt = time.Now()

	return nil
}
