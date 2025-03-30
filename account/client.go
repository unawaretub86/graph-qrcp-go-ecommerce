package account

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/unawaretub86/graph-qrcp-go-ecommerce/account/pb/account"
)

type Client struct {
	conn    *grpc.ClientConn
	Service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		Service: pb.NewAccountServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	res, err := c.Service.GetAccount(ctx, &pb.GetAccountRequest{
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
	r, err := c.Service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}
