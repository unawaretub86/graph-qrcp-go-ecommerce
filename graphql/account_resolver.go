package main

import "context"

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	// accounts := make([]*Account, 0)
	// for _, account := range r.server.accounts {
	// 	accounts = append(accounts, account)
	// }
	// return accounts, nil
	return nil, nil
}