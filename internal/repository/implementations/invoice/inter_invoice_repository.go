package invoice_repository

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/NewLeonardooliv/gateway-payment/internal/domain"
	"github.com/NewLeonardooliv/gateway-payment/internal/shared"
)

type InterInvoiceRepository struct {
	clientID     string
	clientSecret string
	certPath     string
	keyPath      string
	apiUrl       string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewInterInvoiceRepository(clientID, clientSecret string) *InterInvoiceRepository {
	apiUrl := "https://cdpj-sandbox.partners.uatinter.co"
	if shared.GetEnv("ENV", "dev") == "prod" {
		apiUrl = "https://cdpj.partners.bancointer.com.br"
	}

	tlsPath := shared.GetEnv("INTERBANK_TLS_PATH", "cert/")

	certPath := tlsPath + "Sandbox_InterAPI_Certificado.crt"
	keyPath := tlsPath + "Sandbox_InterAPI_Chave.key"

	return &InterInvoiceRepository{clientID, clientSecret, certPath, keyPath, apiUrl}
}

func (r *InterInvoiceRepository) getAccessToken() (*TokenResponse, error) {
	log.Printf("[InterInvoiceRepository] Starting access token request")

	cert, err := tls.LoadX509KeyPair(r.certPath, r.keyPath)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error loading certificate: %v", err)
		return nil, err
	}
	log.Printf("[InterInvoiceRepository] Certificate loaded successfully")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", r.clientID)
	data.Set("client_secret", r.clientSecret)
	data.Set("scope", "cob.write cob.read cobv.write cobv.read lotecobv.write lotecobv.read pix.write pix.read webhook.write webhook.read payloadlocation.write payloadlocation.read boleto-cobranca.read boleto-cobranca.write extrato.read pagamento-pix.write pagamento-pix.read extrato-usend.read pagamento-boleto.read pagamento-boleto.write pagamento-darf.write pagamento-lote.write pagamento-lote.read webhook-banking.read webhook-banking.write")

	log.Printf("[InterInvoiceRepository] Sending token request to %s/oauth/v2/token", r.apiUrl)

	req, err := http.NewRequest("POST", r.apiUrl+"/oauth/v2/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error creating token request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error sending token request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error reading token response: %v", err)
		return nil, err
	}

	log.Printf("[InterInvoiceRepository] Token response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode >= 400 {
		log.Printf("[InterInvoiceRepository] Token request failed with status: %d", resp.StatusCode)
		return nil, fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Printf("[InterInvoiceRepository] Error decoding token response: %v", err)
		return nil, err
	}

	log.Printf("[InterInvoiceRepository] Access token obtained successfully")

	return &tokenResp, nil
}

func (r *InterInvoiceRepository) Save(invoice *domain.Invoice) error {
	log.Printf("[InterInvoiceRepository] Starting boleto creation for invoice: %s", invoice.Reference)

	token, err := r.getAccessToken()
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error obtaining access token: %v", err)

		return err
	}

	log.Printf("[InterInvoiceRepository] Access token obtained successfully")

	payload := map[string]interface{}{
		"seuNumero":      invoice.Reference,
		"valorNominal":   fmt.Sprintf("%.2f", invoice.Amount),
		"dataEmissao":    invoice.CreatedAt.Format("2006-01-02"),
		"dataVencimento": invoice.DueDate.Format("2006-01-02"),
		"pagador": map[string]interface{}{
			"cpfCnpj":    invoice.Payer.TaxID,
			"tipoPessoa": "FISICA",
			"nome":       invoice.Payer.Name,
			"endereco":   invoice.Payer.Address,
			"cidade":     invoice.Payer.City,
			"uf":         invoice.Payer.State,
			"cep":        invoice.Payer.ZipCode,
			"bairro":     invoice.Payer.District,
		},
		"numDiasAgenda": 60,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error serializing payload for invoice %s: %v", invoice.Reference, err)

		return err
	}

	log.Printf("[InterInvoiceRepository] Payload created successfully for invoice %s: %s", invoice.Reference, string(payloadBytes))

	req, err := http.NewRequest("POST", r.apiUrl+"/cobranca/v3/cobrancas", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error creating request for invoice %s: %v", invoice.Reference, err)

		return err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	cert, err := tls.LoadX509KeyPair(r.certPath, r.keyPath)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error loading certificate for invoice %s: %v", invoice.Reference, err)

		return err
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	log.Printf("[InterInvoiceRepository] Sending request to create boleto for invoice %s", invoice.Reference)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error sending request for invoice %s: %v", invoice.Reference, err)

		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[InterInvoiceRepository] Error reading API response for invoice %s: %v", invoice.Reference, err)

		return err
	}

	log.Printf("[InterInvoiceRepository] Inter API response for invoice %s (status %d): %s", invoice.Reference, resp.StatusCode, string(body))

	if resp.StatusCode >= 400 {
		log.Printf("[InterInvoiceRepository] Error creating boleto for invoice %s. Status: %d, Body: %s", invoice.Reference, resp.StatusCode, string(body))

		return fmt.Errorf("error creating Invoice: %s, %d, %s", invoice.Reference, resp.StatusCode, string(body))
	}

	log.Printf("[InterInvoiceRepository] Invoice successfully created for invoice %s", invoice.Reference)

	return nil
}

func (r *InterInvoiceRepository) FindByID(id string) (*domain.Invoice, error) {
	return nil, domain.ErrMethodNotImplemented
}

func (r *InterInvoiceRepository) FindByAccountID(accountID string) ([]*domain.Invoice, error) {
	return nil, domain.ErrMethodNotImplemented
}

func (r *InterInvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	return domain.ErrMethodNotImplemented
}
