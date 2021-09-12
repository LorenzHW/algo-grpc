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

func (s *Server) CreateWallet(ctx context.Context, in *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	walletID, err := s.AlgoInteractor.CreateWallet(in)
	if err != nil {
		return nil, err
	}
	return &pb.CreateWalletResponse{WalletID: walletID}, nil
}

func (s *Server) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	accountAddress, err := s.AlgoInteractor.CreateAccountForWallet(in)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAccountResponse{Address: accountAddress}, nil
}

func (s *Server) MakeTransaction(ctx context.Context, in *pb.MakeTransactionRequest) (*pb.MakeTransactionResponse, error) {
	err := s.AlgoInteractor.MakeTransaction(in)
	if err != nil {
		return nil, err
	}
	return &pb.MakeTransactionResponse{Success: true}, nil
}

func (s *Server) GetTransactions(ctx context.Context, in *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	transactions, err := s.AlgoInteractor.getTransactions(in)
	if err != nil {
		return nil, err
	}
	transactionsMapped := mapToTransactions(transactions.Transactions)
	return &pb.GetTransactionsResponse{Transactions: transactionsMapped}, nil
}

func (s *Server) GetBlock(ctx context.Context, in *pb.GetBlockRequest) (*pb.GetBlockResponse, error) {
	block, err := s.AlgoInteractor.GetBlock(in)
	if err != nil {
		return nil, err
	}
	blockMapped := mapToBlock(block)
	return &pb.GetBlockResponse{Block: blockMapped}, nil
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
