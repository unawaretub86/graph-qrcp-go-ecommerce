package account

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/unawaretub86/graph-qrcp-go-ecommerce/account/pb"
)

type Client struct {
	conn    *grpc.ClientConn
	Server pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		Server: pb.NewAccountServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	res, err := c.Server.GetAccount(ctx, &pb.GetAccountRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   res.Account.Id,
		Name: res.Account.Name,
	}, nil
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.Server.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, take uint64, skip uint64) ([]*Account, error) {
	r, err := c.Server.GetAccounts(ctx, &pb.GetAccountsRequest{Take: take, Skip: skip})
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, &Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}

	return accounts, nil
}
