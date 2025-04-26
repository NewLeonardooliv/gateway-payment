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

type Payer struct {
	ID        string
	Name      string
	TaxID     string
	Email     string
	Phone     string
	Address   string
	Number    string
	District  string
	City      string
	State     string
	ZipCode   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Invoice struct {
	ID             string
	Payer          Payer
	Reference      string
	AccountID      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    string
	CardLastDigits string
	DueDate        time.Time
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

func NewInvoice(accountID string, amount float64, description string, paymentType string, dueDate time.Time, reference string, card CreditCard, payer Payer) (*Invoice, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	cardLastDigits := card.Number[len(card.Number)-4:]

	invoicePayer := &Payer{
		ID:        uuid.New().String(),
		Name:      payer.Name,
		TaxID:     payer.TaxID,
		Email:     payer.Email,
		Phone:     payer.Phone,
		Address:   payer.Address,
		Number:    payer.Number,
		District:  payer.District,
		City:      payer.City,
		State:     payer.State,
		ZipCode:   payer.ZipCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Amount:         amount,
		Status:         "pending",
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: cardLastDigits,
		Payer:          *invoicePayer,
		DueDate:        dueDate,
		Reference:      reference,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (invoice *Invoice) Process() error {
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
