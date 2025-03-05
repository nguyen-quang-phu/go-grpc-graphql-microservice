package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var ErrNotFound = errors.New("Entity not found")

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elasticsearch.TypedClient
}

type productDocument struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
	})
	if err != nil {
		return nil, err
	}

	return &elasticRepository{client}, nil
}

func (repository *elasticRepository) Close() {}

func (repository *elasticRepository) PutProduct(ctx context.Context, product Product) error {
	data, err := json.Marshal(productDocument{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})
	if err != nil {
		log.Fatalf("Error marshaling product data: %s", err)
	}

	_, err = repository.client.Index("products").Id(product.ID).Request(data).Do(ctx)
	return err
}

func (repository *elasticRepository) GetProductById(
	ctx context.Context,
	id string,
) (*Product, error) {
	res, err := repository.client.Get("products", id).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	product := productDocument{}
	if err = json.Unmarshal(*&res.Source_, &product); err != nil {
		return nil, err
	}

	return &Product{
		ID:    id,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

func (repository *elasticRepository) ListProducts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Product, error) {
	res, err := repository.client.Search().
		Index("products").
		Request(&search.Request{Query: &types.Query{MatchAll: &types.MatchAllQuery{}}}).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*&hit.Source_, &p); err == nil {
			products = append(products, Product{
				ID:          *hit.Id_,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}

	return products, nil
}

func (repository *elasticRepository) ListProductsWithIDs(
	ctx context.Context,
	ids []string,
) ([]Product, error) {
	res, err := repository.client.Search().
		Index("products").
		Request(&search.Request{
			Query: &types.Query{
				Ids: &types.IdsQuery{
					Values: ids,
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*&hit.Source_, &p); err == nil {
			products = append(products, Product{
				ID:          *hit.Id_,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}

	return products, nil
}

func (repository *elasticRepository) SearchProducts(
	ctx context.Context,
	query string,
	skip uint64,
	take uint64,
) ([]Product, error) {
	res, err := repository.client.Search().
		Index("products").
		Request(&search.Request{
			Query: &types.Query{
				MultiMatch: &types.MultiMatchQuery{
					Query:  query,
					Fields: []string{"name", "description"},
				},
			},
		}).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*&hit.Source_, &p); err == nil {
			products = append(products, Product{
				ID:          *hit.Id_,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}

	return products, nil
}
