package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	ordersList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	for _, order := range ordersList {
		var products []*OrderedProduct
		for _, product := range order.Products {
			products = append(products, &OrderedProduct{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    int(product.Quantity),
			})
		}

		orders = append(orders, &Order{
			ID:         order.ID,
			Products:   products,
			TotalPrice: order.TotalPrice,
			CreatedAt:  order.CreatedAt,
		})
	}

	return orders, nil
}
