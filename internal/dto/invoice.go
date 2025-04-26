package dto

import (
	"time"

	"github.com/NewLeonardooliv/gateway-payment/internal/domain"
)

const (
	StatusPending  = string(domain.StatusPending)
	StatusApproved = string(domain.StatusApproved)
	StatusRejected = string(domain.StatusRejected)
)

type CreateInvoiceInput struct {
	APIKey         string
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
	PaymentType    string    `json:"payment_type"`
	CardNumber     string    `json:"card_number"`
	CVV            string    `json:"cvv"`
	ExpiryMonth    int       `json:"expiry_month"`
	ExpiryYear     int       `json:"expiry_year"`
	CardholderName string    `json:"cardholder_name"`
	DueDate        time.Time `json:"due_date"`
	Reference      string    `json:"reference"`
	Payer          struct {
		Name     string `json:"name"`
		TaxID    string `json:"tax_id"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
		Number   string `json:"number"`
		District string `json:"district"`
		City     string `json:"city"`
		State    string `json:"state"`
		ZipCode  string `json:"zip_code"`
	} `json:"payer"`
}

type PayerOutput struct {
	Name     string `json:"name"`
	TaxID    string `json:"tax_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Number   string `json:"number"`
	District string `json:"district"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
}

type InvoiceOutput struct {
	ID             string      `json:"id"`
	AccountID      string      `json:"account_id"`
	Amount         float64     `json:"amount"`
	Status         string      `json:"status"`
	Description    string      `json:"description"`
	PaymentType    string      `json:"payment_type"`
	CardLastDigits string      `json:"card_last_digits"`
	Reference      string      `json:"reference"`
	DueDate        time.Time   `json:"due_date"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      *time.Time  `json:"deleted_at"`
	Payer          PayerOutput `json:"payer"`
}

func ToInvoice(input CreateInvoiceInput, accountID string) (*domain.Invoice, error) {
	card := domain.CreditCard{
		Number:         input.CardNumber,
		CVV:            input.CVV,
		ExpiryMonth:    input.ExpiryMonth,
		ExpiryYear:     input.ExpiryYear,
		CardholderName: input.CardholderName,
	}

	payer := domain.Payer{
		Name:     input.Payer.Name,
		TaxID:    input.Payer.TaxID,
		Email:    input.Payer.Email,
		Phone:    input.Payer.Phone,
		Address:  input.Payer.Address,
		Number:   input.Payer.Number,
		District: input.Payer.District,
		City:     input.Payer.City,
		State:    input.Payer.State,
		ZipCode:  input.Payer.ZipCode,
	}

	return domain.NewInvoice(
		accountID,
		input.Amount,
		input.Description,
		input.PaymentType,
		input.DueDate,
		input.Reference,
		card,
		payer,
	)
}

func FromInvoice(invoice *domain.Invoice) *InvoiceOutput {
	Payer := PayerOutput{
		Name:     invoice.Payer.Name,
		TaxID:    invoice.Payer.TaxID,
		Email:    invoice.Payer.Email,
		Phone:    invoice.Payer.Phone,
		Address:  invoice.Payer.Address,
		Number:   invoice.Payer.Number,
		District: invoice.Payer.District,
		City:     invoice.Payer.City,
		State:    invoice.Payer.State,
		ZipCode:  invoice.Payer.ZipCode,
	}

	return &InvoiceOutput{
		ID:             invoice.ID,
		AccountID:      invoice.AccountID,
		Amount:         invoice.Amount,
		Status:         string(invoice.Status),
		Description:    invoice.Description,
		PaymentType:    invoice.PaymentType,
		CardLastDigits: invoice.CardLastDigits,
		DueDate:        invoice.DueDate,
		Payer:          Payer,
		Reference:      invoice.Reference,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}
}
