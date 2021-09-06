package main

import (
	"context"
	pb "github.com/LorenzHW/algo-grpc/protos"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedAlgorandServer
	AlgoInteractor *AlgoInteractor
}

func NewServer() (server *Server) {
	server = &Server{}
	server.AlgoInteractor = NewAlgoInteractor()
	return server
}

func (s *Server) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.AlgoInteractor.GetAccount(in)
	if err != nil {
		return nil, err
	}
	accountMapped := mapToAccount(account)
	return &pb.GetAccountResponse{Account: accountMapped}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := NewServer()

	pb.RegisterAlgorandServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
