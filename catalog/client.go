package catalog

import (
	"google.golang.org/grpc"

	pb "github.com/unawaretub86/graph-qrcp-go-ecommerce/catalog/pb"
)

type Client struct {
	conn    *grpc.ClientConn
	Service pb.CatalogClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		Service: pb.NewCatalogClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}