package catalog

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/unawaretub86/graph-qrcp-go-ecommerce/catalog/pb"
)

type grpcServer struct {
	pb.UnimplementedCatalogServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterCatalogServer(srv, &grpcServer{
		UnimplementedCatalogServer: pb.UnimplementedCatalogServer{},
		service:                    s,
	})
	reflection.Register(srv)

	return srv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p := &Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}

	a, err := s.service.PostProduct(ctx, p)
	if err != nil {
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          a.ID,
			Name:        a.Name,
			Price:       a.Price,
			Description: a.Description,
		},
	}, nil
}

func (s *grpcServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	a, err := s.service.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductByIDResponse{
		Product: &pb.Product{
			Id:          a.ID,
			Name:        a.Name,
			Price:       a.Price,
			Description: a.Description,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []*Product
	var err error

	if req.Query != "" {
		res, err = s.service.SearchProducts(ctx, req.Query, req.Take, req.Skip)
		if err != nil {
			return nil, err
		}
	} else if len(req.Ids) > 0 {
		res, err = s.service.GetProductsByIDs(ctx, req.Ids)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = s.service.GetProducts(ctx, req.Take, req.Skip)
		if err != nil {
			return nil, err
		}
	}

	pbProducts := []*pb.Product{}
	for _, p := range res {
		pbProducts = append(pbProducts, &pb.Product{
			Id:   p.ID,
			Name: p.Name,
		})
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
