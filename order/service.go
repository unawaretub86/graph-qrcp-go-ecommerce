package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(context.Context, string, []*OrderedProduct) (*Order, error)
	GetOrdersForAccount(context.Context, string) ([]*Order, error)
}

type Order struct {
	ID         string            `json:"id"`
	CreatedAt  time.Time         `json:"created_at"`
	TotalPrice float64           `json:"total_price"`
	AccountID  string            `json:"account_id"`
	Products   []*OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
}

type orderService struct {
	repository Repository
}

func (o *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	return o.repository.GetOrdersForAccount(ctx, accountId)
}

func (o *orderService) PostOrder(ctx context.Context, accountId string, products []*OrderedProduct) (*Order, error) {
	order := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now(),
		AccountID: accountId,
		Products:  products,
	}

	order.TotalPrice = 0
	for _, product := range products {
		order.TotalPrice += product.Price * float64(product.Quantity)
	}

	err := o.repository.PutOrder(ctx, *order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func NewService(repository Repository) Service {
	return &orderService{repository}
}
