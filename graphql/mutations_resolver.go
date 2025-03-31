package main

import (
	"context"
	"errors"
	"time"

	"github.com/unawaretub86/graph-qrcp-go-ecommerce/order"
)

var errorNotFound = errors.New("not found")

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input *AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	account, err := r.server.accountClient.PostAccount(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input *OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	var products []*order.OrderedProduct
	for _, product := range input.Products {
		if product.Quantity <= 0 {
			return nil, errorNotFound
		}

		products = append(products, &order.OrderedProduct{
			ID:       product.ID,
			Quantity: int64(product.Quantity),
		})
	}

	orderResult, err := r.server.orderClient.PostOrder(ctx, input.AccountID, products)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         orderResult.ID,
		CreatedAt:  orderResult.CreatedAt,
		TotalPrice: orderResult.TotalPrice,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input *ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	product, err := r.server.catalogClient.PostProduct(ctx, input.Name, *input.Description, input.Price)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: &product.Description,
	}, nil
}
