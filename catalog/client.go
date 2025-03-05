package catalog

import (
	"context"

	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewProductServiceClient(conn)
	return &Client{conn, service}, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) PostProduct(
	ctx context.Context,
	name, description string,
	price float64,
) (*Product, error) {
	response, err := client.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          response.Product.Id,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
	}, nil
}

func (client *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	response, err := client.service.GetProduct(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          response.Product.Id,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
	}, nil
}

func (client *Client) GetProducts(
	ctx context.Context,
	skip, take uint64,
	ids []string,
	query string,
) ([]Product, error) {
	response, err := client.service.GetProducts(ctx, &pb.GetProductsRequest{
		Ids:   ids,
		Skip:  skip,
		Take:  take,
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, product := range response.Products {
		products = append(products, Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, nil
}
