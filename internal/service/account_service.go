package service

import "github.com/NewLeonardooliv/gateway-payment/internal/domain"

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{
		repository: repository,
	}
}
