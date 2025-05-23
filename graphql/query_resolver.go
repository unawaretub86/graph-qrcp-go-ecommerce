package main

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Account{{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for _, account := range accountList {
		accounts = append(accounts, &Account{
			ID:   account.ID,
			Name: account.Name,
		})
	}

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: &r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	var products []*Product
	for _, product := range productList {
		products = append(products, &Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: &product.Description,
			Price:       product.Price,
		})
	}

	return products, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(0)

	if p.Skip > 0 {
		skipValue = uint64(p.Skip)
	}

	if p.Take > 0 {
		takeValue = uint64(p.Take)
	}

	return skipValue, takeValue
}
