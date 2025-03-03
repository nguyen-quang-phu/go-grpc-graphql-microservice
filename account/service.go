package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccountById(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type accountService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &accountService{repository}
}

func (service *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	if err := service.repository.PutAccount(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (service *accountService) GetAccountById(ctx context.Context, id string) (*Account, error) {
	return service.repository.GetAccountById(ctx, id)
}

func (service *accountService) GetAccounts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return service.repository.ListAccounts(ctx, skip, take)
}
