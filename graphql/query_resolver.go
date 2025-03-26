package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	// accounts := make([]*Account, 0)
	// for _, account := range r.server.accounts {
	// 	accounts = append(accounts, account)
	// }
	// return accounts, nil
	return nil, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	// accounts := make([]*Account, 0)
	// for _, account := range r.server.accounts {
	// 	accounts = append(accounts, account)
	// }
	// return accounts, nil
	return nil, nil
}

