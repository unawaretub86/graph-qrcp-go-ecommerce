package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	PutAccount(context.Context, *Account) error
	GetAccountById(context.Context, string) (*Account, error)
	ListAccounts(context.Context, uint64, uint64) ([]*Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) PutAccount(ctx context.Context, account *Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES ($1, $2)", account.ID, account.Name)
	return err
}

func (r *postgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	if err := rows.Scan(&account.ID, &account.Name); err != nil {
		return nil, err
	}

	return account, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id LIMIT $1 OFFSET $2", take, skip)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func NewPostgresRepository(dbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}
