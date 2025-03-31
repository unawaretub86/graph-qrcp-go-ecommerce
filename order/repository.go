package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(context.Context, Order) error
	GetOrdersForAccount(context.Context, string) ([]*Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) PutOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		order.ID,
		order.CreatedAt,
		order.AccountID,
		order.TotalPrice,
	)
	if err != nil {
		return err
	}

	stm, err := tx.PrepareContext(
		ctx,
		pq.CopyIn(
			"order_products",
			"order_id",
			"product_id",
			"quantity",
		),
	)
	if err != nil {
		return err
	}

	for _, product := range order.Products {
		_, err = stm.ExecContext(
			ctx,
			order.ID,
			product.ID,
			product.Quantity,
		)
		if err != nil {
			return err
		}
	}

	_, err = stm.ExecContext(ctx)
	if err != nil {
		return err
	}

	err = stm.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id,
		o.account_id,
		o.created_at,
		o.total_price::money::numeric::float8,
		op.quantity
		op.product_id
		FROM orders as o 
		JOIN order_products as op
		ON (o.id = op.order_id)
		WHERE o.account_id = $1
		ORDER BY o.id`, accountID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []*Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []*OrderedProduct{}

	for rows.Next() {
		if err := rows.Scan(
			&order.ID,
			&order.AccountID,
			&order.CreatedAt,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != "" && order.ID != lastOrder.ID {
			newOrder := &Order{
				ID:         order.ID,
				CreatedAt:  order.CreatedAt,
				AccountID:  order.AccountID,
				TotalPrice: order.TotalPrice,
				Products:   products,
			}

			orders = append(orders, newOrder)
			products = []*OrderedProduct{}
		}

		products = append(products, &OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		*lastOrder = *order
	}

	newOrder := &Order{
		ID:         lastOrder.ID,
		CreatedAt:  lastOrder.CreatedAt,
		AccountID:  lastOrder.AccountID,
		TotalPrice: lastOrder.TotalPrice,
		Products:   lastOrder.Products,
	}

	orders = append(orders, newOrder)

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
