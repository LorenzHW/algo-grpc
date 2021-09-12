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

	r, err := c.GetAccount(ctx, &pb.GetAccountRequest{AccountAddress: "KU2PUHUOFNTUIEZKIIVHRVRXAH6EISEBA2IZRLQQQOD3626PBVVMW3TCCU"})
	if err != nil {
		log.Fatalf("could not get account info: %v", err)
	}
	printResponse(r)
}

func printResponse(response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("could not convert to json: %v", err)
	}
	log.Printf(string(b))
}
