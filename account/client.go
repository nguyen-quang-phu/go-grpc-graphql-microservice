package account

import (
	"context"

	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewAccountServiceClient(conn)
	return &Client{conn, service}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	response, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return &Account{ID: response.Account.Id, Name: response.Account.Name}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	response, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &Account{ID: response.Account.Id, Name: response.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	response, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0, len(response.Accounts))
	for _, a := range response.Accounts {
		accounts = append(accounts, &Account{ID: a.Id, Name: a.Name})
	}

	return accounts, nil
}
