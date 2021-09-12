package main

import (
	"context"
	"encoding/json"
	pb "github.com/LorenzHW/algo-grpc/protos"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAlgorandClient(conn)

	clientDeadline := time.Now().Add(time.Duration(10000000) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	// Comment/uncomment functions as you wish
	getAccount(c, ctx)
	//createAccount(c, ctx)
	//makeTransaction(c, ctx)
	//getTransaction(c, ctx)
	//getBlock(c, ctx)
}

func createAccount(c pb.AlgorandClient, ctx context.Context) {
	walletPassword := "dummyPassword"
	walletID := createWallet(c, ctx, walletPassword)

	r, err := c.CreateAccount(ctx, &pb.CreateAccountRequest{
		WalletID:       walletID,
		WalletPassword: walletPassword,
	})
	if err != nil {
		log.Fatalf("Could not create account : %v", err)
		return
	}
	printResponse(r)
}

func makeTransaction(c pb.AlgorandClient, ctx context.Context) {
	log.Printf("Making a transaction")
	makeTransactionResponse, err := c.MakeTransaction(ctx, &pb.MakeTransactionRequest{FromWalletAddress: "JDXV7FAW5Y7TU7VDSPU5AP3NPTJEBCMTL2VUEMYZWON5JDEMCHMGACXNTQ", FromWalletID: "1bb14a2d8abb3380e8a5dbaccbd17773", FromWalletPassword: "", ToWalletAddress: "KU2PUHUOFNTUIEZKIIVHRVRXAH6EISEBA2IZRLQQQOD3626PBVVMW3TCCU"})
	if err != nil {
		log.Fatalf("Could not make transaction: %v", err)
		return
	}
	printResponse(makeTransactionResponse)
}

func getAccount(c pb.AlgorandClient, ctx context.Context) {
	log.Printf("Getting an account")
	r, err := c.GetAccount(ctx, &pb.GetAccountRequest{AccountAddress: "KU2PUHUOFNTUIEZKIIVHRVRXAH6EISEBA2IZRLQQQOD3626PBVVMW3TCCU"})
	if err != nil {
		log.Fatalf("could not get account info: %v", err)
	}
	printResponse(r)
}

func createWallet(c pb.AlgorandClient, ctx context.Context, password string) string {
	log.Printf("Creating a wallet")
	r, err := c.CreateWallet(ctx, &pb.CreateWalletRequest{WalletName: "MyAlgoWallet2", WalletPassword: password})
	if err != nil {
		log.Fatalf("could not create wallet: %v", err)
		return ""
	}
	log.Printf("The Wallet ID is: %s", r.GetWalletID())
	return r.GetWalletID()
}

func getTransaction(c pb.AlgorandClient, ctx context.Context) {
	log.Printf("Getting a transaction")
	transactionResponse, err := c.GetTransactions(ctx, &pb.GetTransactionsRequest{Address: "JDXV7FAW5Y7TU7VDSPU5AP3NPTJEBCMTL2VUEMYZWON5JDEMCHMGACXNTQ"})
	if err != nil {
		log.Fatalf("could not get transaction history: %v", err)
		return
	}
	printResponse(transactionResponse)
}

func getBlock(c pb.AlgorandClient, ctx context.Context) {
	log.Printf("Getting a block")
	blockResponse, err := c.GetBlock(ctx, &pb.GetBlockRequest{RoundNumber: 48})
	if err != nil {
		log.Fatalf("could not get transaction history: %v", err)
		return
	}
	printResponse(blockResponse)
}

func printResponse(response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("could not convert to json: %v", err)
		return
	}
	log.Printf(string(b))
}
