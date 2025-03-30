package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(context.Context, *Product) (*Product, error)
	GetProductByID(context.Context, string) (*Product, error)
	GetProducts(context.Context, uint64, uint64) ([]*Product, error)
	GetProductsByIDs(context.Context, []string) ([]*Product, error)
	SearchProducts(context.Context, string, uint64, uint64) ([]*Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type catalogService struct {
	repository Repository
}

func (c *catalogService) GetProductByID(ctx context.Context, id string) (*Product, error) {
	return c.repository.GetProductByID(ctx, id)
}

func (c *catalogService) GetProducts(ctx context.Context, take uint64, skip uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return c.repository.ListProducts(ctx, skip, take)
}

func (c *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]*Product, error) {
	return c.repository.ListProductsWithIDs(ctx, ids)
}

func (c *catalogService) PostProduct(ctx context.Context, product *Product) (*Product, error) {
	product.ID = ksuid.New().String()

	if err := c.repository.PutProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (c *catalogService) SearchProducts(ctx context.Context, query string, take uint64, skip uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return c.repository.SearchProducts(ctx, query, skip, take)
}

func NewService(repository Repository) Service {
	return &catalogService{repository}
}
