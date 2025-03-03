package account

import (
	"context"
	"fmt"
	"net"

	"github.com/nguyen-quang-phu/go-grpc-graphql-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(service Service, port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		service:                           service,
	})
	reflection.Register(server)
	return server.Serve(listen)
}

func (server *grpcServer) PostAccount(
	ctx context.Context,
	request *pb.PostAccountRequest,
) (*pb.PostAccountResponse, error) {
	account, err := server.service.PostAccount(ctx, request.Name)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{Account: &pb.Account{
		Id:   account.ID,
		Name: account.Name,
	}}, nil
}

func (server *grpcServer) GetAccount(
	ctx context.Context,
	request *pb.GetAccountRequest,
) (*pb.GetAccountResponse, error) {
	account, err := server.service.GetAccountById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{Account: &pb.Account{
		Id:   account.ID,
		Name: account.Name,
	}}, nil
}

func (server *grpcServer) GetAccounts(
	ctx context.Context,
	request *pb.GetAccountsRequest,
) (*pb.GetAccountsResponse, error) {
	accs, err := server.service.GetAccounts(ctx, request.Skip, request.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, acc := range accs {
		accounts = append(accounts, &pb.Account{
			Id:   acc.ID,
			Name: acc.Name,
		})
	}

	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
