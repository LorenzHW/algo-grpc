package main

import (
	"context"
	"fmt"
	pb "github.com/LorenzHW/algo-grpc/protos"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
)

type AlgoInteractor struct {
	kmdClient     kmd.Client
	algodClient   algod.Client
	indexerClient indexer.Client
}

func NewAlgoInteractor() (algoInteractor *AlgoInteractor) {
	algoInteractor = &AlgoInteractor{}
	algodAddress := "http://localhost:4001"
	kmdAddress := "http://localhost:4002"
	indexerAddress := "http://localhost:8980"

	algodToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	kmdToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	fmt.Println("Made a kmd client")
	algoInteractor.kmdClient = kmdClient

	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}
	fmt.Println("Made an algodClient")
	algoInteractor.algodClient = *algodClient

	algoIndexer, err := indexer.MakeClient(indexerAddress, "")
	if err != nil {
		return
	}
	fmt.Println("Made an indexerClient")
	algoInteractor.indexerClient = *algoIndexer

	return algoInteractor
}

func (algoInteractor *AlgoInteractor) GetAccount(getUserAccountRequest *pb.GetAccountRequest) (models.Account, error) {
	return algoInteractor.algodClient.AccountInformation(getUserAccountRequest.WalletAddress).Do(context.Background())
}
