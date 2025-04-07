package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func generateAPIKey() string {
	b := make([]byte, 16)
	rand.Read(b)

	return hex.EncodeToString(b)
}

func NewAccount(name, email string) *Account {
	account := &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		APIKey:    generateAPIKey(),
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

func (account *Account) AddBalance(amount float64) {
	account.mu.Lock()
	defer account.mu.Unlock()

	account.Balance += amount
	account.UpdatedAt = time.Now()
}
