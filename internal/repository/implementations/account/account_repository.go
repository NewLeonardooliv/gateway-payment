package account_repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/NewLeonardooliv/gateway-payment/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	accountRepository := &AccountRepository{
		db: db,
	}

	return accountRepository
}

func (repository *AccountRepository) Save(account *domain.Account) error {
	statement, err := repository.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6 ,$7)
	`)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	log.Printf("Finding account by API key: %s", apiKey)

	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE api_key = $1
			AND deleted_at IS NULL
	`, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("Account not found with API key: %s", apiKey)
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		log.Printf("Error finding account with API key %s: %v", apiKey, err)
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	log.Printf("Account found: ID=%s, Name=%s, Email=%s, APIKey=%s, Balance=%.2f",
		account.ID, account.Name, account.Email, account.APIKey, account.Balance)
	return &account, nil
}

func (repository *AccountRepository) FindByID(id string) (*domain.Account, error) {
	log.Printf("Finding account by ID: %s", id)

	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
			AND deleted_at IS NULL
	`, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("Account not found with ID: %s", id)
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		log.Printf("Error finding account with ID %s: %v", id, err)
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	log.Printf("Account found: ID=%s, Name=%s, Email=%s, APIKey=%s, Balance=%.2f",
		account.ID, account.Name, account.Email, account.APIKey, account.Balance)
	return &account, nil
}

func (repository *AccountRepository) UpdateBalance(account *domain.Account) error {
	log.Printf("Updating balance for account ID %s to %.2f", account.ID, account.Balance)

	tx, err := repository.db.Begin()

	if err != nil {
		log.Printf("Error starting transaction for account %s: %v", account.ID, err)
		return err
	}

	defer tx.Rollback()

	var currentBalance float64

	err = tx.QueryRow(`
		SELECT balance
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`, account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		log.Printf("Account not found for balance update: %s", account.ID)
		return domain.ErrAccountNotFound
	}

	if err != nil {
		log.Printf("Error querying current balance for account %s: %v", account.ID, err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3
	`, account.Balance, time.Now(), account.ID)

	if err != nil {
		log.Printf("Error updating balance for account %s: %v", account.ID, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction for account %s: %v", account.ID, err)
		return err
	}

	log.Printf("Balance updated successfully for account %s", account.ID)
	return nil
}
