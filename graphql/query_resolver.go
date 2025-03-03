package main

import "context"

type queryResolver struct {
	server *Server
}

func (resolver *queryResolver) Accounts(
	ctx context.Context,
	pagination *PagninationInput,
	id *string,
) ([]*Account, error) {
	return nil, nil
}

func (resolver *queryResolver) Products(
	ctx context.Context,
	pagination *PagninationInput,
	query *string,
	id *string,
) ([]*Product, error) {
	return nil, nil
}
