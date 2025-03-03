package main

import "context"

type mutationResolver struct {
	server *Server
}

func (resolver *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
	
	return nil,nil
}


func (resolver *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {

	return nil,nil
}


func (resolver *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {

	return nil,nil
}
