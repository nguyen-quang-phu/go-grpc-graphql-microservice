package catalog

import (
	"context"
	"fmt"
	"net"

	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenGRPC(service Service, port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, &grpcServer{
		UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{},
		service:                           service,
	})
	reflection.Register(server)
	return server.Serve(listen)
}

func (server *grpcServer) PostProduct(
	ctx context.Context,
	req *pb.PostProductRequest,
) (*pb.PostProductResponse, error) {
	product, err := server.service.PostProduct(ctx, req.Name, req.Description, float64(req.Price))
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil
}

func (server *grpcServer) GetProduct(
	ctx context.Context,
	req *pb.GetProductRequest,
) (*pb.GetProductResponse, error) {
	product, err := server.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil
}

func (server *grpcServer) handleProductRequest(
	ctx context.Context,
	req *pb.GetProductsRequest,
) ([]Product, error) {
	if req.Query != "" {
		return server.service.SearchProducts(ctx, req.Query, req.Take, req.Skip)
	}

	if len(req.Ids) != 0 {
		return server.service.GetProductsByIDs(ctx, req.Ids)
	}

	return server.service.GetProducts(ctx, req.Take, req.Skip)
}

func (server *grpcServer) GetProducts(
	ctx context.Context,
	req *pb.GetProductsRequest,
) (*pb.GetProductsResponse, error) {
	res, err := server.handleProductRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}
	for _, product := range res {
		products = append(products, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return &pb.GetProductsResponse{Products: products}, nil
}
