package domain

type PaymentMethod string

const (
	PaymentMethodBoleto PaymentMethod = "boleto"
	PaymentMethodCard   PaymentMethod = "card"
)

type PaymentRequest struct {
	Amount      float64
	Currency    string
	Description string
	Method      PaymentMethod
	Customer    CustomerInfo
	Metadata    map[string]string
}

type CustomerInfo struct {
	Name      string
	Email     string
	CPFOrCNPJ string
}

type PaymentResponse struct {
	ID          string
	Status      string
	RedirectURL string
	PDFURL      string
}

type PaymentProvider interface {
	CreatePayment(req PaymentRequest) (*PaymentResponse, error)
}
