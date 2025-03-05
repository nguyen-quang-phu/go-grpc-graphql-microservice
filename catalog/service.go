package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, take, skip uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, take, skip uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

type catalogService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &catalogService{repository}
}

func (service *catalogService) PostProduct(
	ctx context.Context,
	name, description string,
	price float64,
) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
	}
	if err := service.repository.PutProduct(ctx, *product); err != nil {
		return nil, err
	}

	return product, nil
}

func (service *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return service.repository.GetProductById(ctx, id)
}

func (service *catalogService) GetProducts(
	ctx context.Context,
	take, skip uint64,
) ([]Product, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}
	return service.repository.ListProducts(ctx, skip, take)
}

func (service *catalogService) GetProductsByIDs(
	ctx context.Context,
	ids []string,
) ([]Product, error) {
	return service.repository.ListProductsWithIDs(ctx, ids)
}

func (service *catalogService) SearchProducts(
	ctx context.Context,
	query string,
	skip, take uint64,
) ([]Product, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	return service.repository.SearchProducts(ctx, query, skip, take)
}
