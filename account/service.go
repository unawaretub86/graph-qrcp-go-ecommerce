package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(context.Context, string) (*Account, error)
	GetAccount(context.Context, string) (*Account, error)
	GetAccounts(context.Context, uint64, uint64) ([]*Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccountService struct {
	repository Repository
}

func (a *AccountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	account, err := a.repository.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountService) GetAccounts(ctx context.Context, take uint64, skip uint64) ([]*Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return a.repository.ListAccounts(ctx, take, skip)
}

func (a *AccountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		ID:   ksuid.New().String(),
		Name: name,
	}

	if err := a.repository.PutAccount(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func NewAccountService(r Repository) Service {
	return &AccountService{r}
}
