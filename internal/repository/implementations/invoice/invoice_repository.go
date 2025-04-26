package invoice_repository

import (
	"database/sql"
	"log"

	"github.com/NewLeonardooliv/gateway-payment/internal/domain"
)

type PostgresInvoiceRepository struct {
	db *sql.DB
}

func NewPostgresInvoiceRepository(db *sql.DB) *PostgresInvoiceRepository {
	return &PostgresInvoiceRepository{
		db: db,
	}
}

func (repository *PostgresInvoiceRepository) Save(invoice *domain.Invoice) error {
	log.Printf("Saving invoice: %+v", invoice)

	_, err := repository.db.Exec(
		"INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, due_date, reference, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		invoice.ID,
		invoice.AccountID,
		invoice.Amount,
		invoice.Status,
		invoice.Description,
		invoice.PaymentType,
		invoice.CardLastDigits,
		invoice.DueDate,
		invoice.Reference,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error saving invoice %s: %v", invoice.ID, err)

		return err
	}

	log.Printf("Invoice saved successfully: %s", invoice.ID)

	_, err = repository.db.Exec(
		"INSERT INTO payers (id, invoice_id, name, tax_id, email, phone, address, number, district, city, state, zip_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
		invoice.Payer.ID,
		invoice.ID,
		invoice.Payer.Name,
		invoice.Payer.TaxID,
		invoice.Payer.Email,
		invoice.Payer.Phone,
		invoice.Payer.Address,
		invoice.Payer.Number,
		invoice.Payer.District,
		invoice.Payer.City,
		invoice.Payer.State,
		invoice.Payer.ZipCode,
		invoice.Payer.CreatedAt,
		invoice.Payer.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error saving payer %s: %v", invoice.Payer.ID, err)
		return err
	}

	log.Printf("Payer saved successfully: %s", invoice.Payer.ID)

	return nil
}

func (r *PostgresInvoiceRepository) FindByID(id string) (*domain.Invoice, error) {
	log.Printf("Finding invoice by ID: %s", id)

	var invoice domain.Invoice
	err := r.db.QueryRow(`
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`, id).Scan(
		&invoice.ID,
		&invoice.AccountID,
		&invoice.Amount,
		&invoice.Status,
		&invoice.Description,
		&invoice.PaymentType,
		&invoice.CardLastDigits,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("Invoice not found: %s", id)
		return nil, domain.ErrInvoiceNotFound
	}

	if err != nil {
		log.Printf("Error finding invoice %s: %v", id, err)

		return nil, err
	}

	log.Printf("Invoice found: %+v", invoice)

	return &invoice, nil
}

func (r *PostgresInvoiceRepository) FindByAccountID(accountID string) ([]*domain.Invoice, error) {
	log.Printf("FindByAccountID called with accountID: %s", accountID)

	rows, err := r.db.Query(`
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE account_id = $1
	`, accountID)

	if err != nil {
		log.Printf("Error executing query for accountID %s: %v", accountID, err)

		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var invoices []*domain.Invoice
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.AccountID, &invoice.Amount, &invoice.Status, &invoice.Description, &invoice.PaymentType, &invoice.CardLastDigits, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row for accountID %s: %v", accountID, err)

			return nil, err
		}

		log.Printf("Fetched invoice: %+v", invoice)
		invoices = append(invoices, &invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error for accountID %s: %v", accountID, err)

		return nil, err
	}

	log.Printf("FindByAccountID completed successfully for accountID %s, total invoices found: %d", accountID, len(invoices))

	return invoices, nil
}

func (r *PostgresInvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	log.Printf("Updating status for invoice ID %s to %s", invoice.ID, invoice.Status)

	result, err := r.db.Exec(
		"UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3",
		invoice.Status, invoice.UpdatedAt, invoice.ID,
	)
	if err != nil {
		log.Printf("Error updating status for invoice %s: %v", invoice.ID, err)

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for invoice %s: %v", invoice.ID, err)

		return err
	}

	if rowsAffected == 0 {
		log.Printf("Invoice not found for update: %s", invoice.ID)

		return domain.ErrInvoiceNotFound
	}

	log.Printf("Status updated successfully for invoice %s", invoice.ID)

	return nil
}
