package main

import (
	"context"
	"fmt"
	pb "github.com/LorenzHW/algo-grpc/protos"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
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
	return algoInteractor.algodClient.AccountInformation(getUserAccountRequest.AccountAddress).Do(context.Background())
}

func (algoInteractor *AlgoInteractor) CreateWallet(createWalletRequest *pb.CreateWalletRequest) (string, error) {

	// Create the example wallet, if it doesn't already exist
	cwResponse, err := algoInteractor.kmdClient.CreateWallet(createWalletRequest.WalletName, createWalletRequest.WalletPassword, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return "", err
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	exampleWalletID := cwResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, exampleWalletID)
	return exampleWalletID, nil
}

func (algoInteractor *AlgoInteractor) CreateAccountForWallet(createAccountRequest *pb.CreateAccountRequest) (string, error) {
	walletHandleToken, err := algoInteractor.getWalletHandleToken(createAccountRequest.WalletID, createAccountRequest.WalletPassword)
	if err != nil {
		return "", err
	}

	// Generate a new address from the wallet handle
	genResponse, err := algoInteractor.kmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return "", err
	}
	fmt.Printf("Generated address %s\n", genResponse.Address)
	return genResponse.Address, nil
}

func (algoInteractor *AlgoInteractor) MakeTransaction(makeTransactionRequest *pb.MakeTransactionRequest) (error error) {
	walletHandleToken, err := algoInteractor.getWalletHandleToken(makeTransactionRequest.FromWalletID, makeTransactionRequest.FromWalletPassword)
	if err != nil {
		fmt.Printf("Error getting wallet handle token: %s\n", err)
		return err
	}

	// Get the suggested transaction parameters
	txParams, err := algoInteractor.algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error building suggested params: %s\n", err)
		return err
	}

	tx, err := future.MakePaymentTxn(makeTransactionRequest.FromWalletAddress, makeTransactionRequest.ToWalletAddress, 150000, nil, "", txParams)
	if err != nil {
		fmt.Printf("Error making payment txn: %s\n", err)
		return err
	}

	signResponse, err := algoInteractor.kmdClient.SignTransaction(walletHandleToken, makeTransactionRequest.FromWalletPassword, tx)
	if err != nil {
		fmt.Printf("Error signing txn: %s\n", err)
		return err
	}
	fmt.Printf("kmd made signed transaction with bytes: %x\n", signResponse.SignedTransaction)

	sendResponse, err := algoInteractor.algodClient.SendRawTransaction(signResponse.SignedTransaction).Do(context.Background())
	if err != nil {
		fmt.Printf("Error sending txn to network: %s\n", err)
		return err
	}
	fmt.Printf("Transaction ID: %s\n", sendResponse)
	return nil
}

func (algoInteractor *AlgoInteractor) GetBlock(in *pb.GetBlockRequest) (models.Block, error) {
	return algoInteractor.indexerClient.LookupBlock(in.RoundNumber).Do(context.Background())
}

// getWalletHandleToken gets a wallet handle. The wallet handle is used for things like signing transactions
// and creating accounts. Wallet handles do expire, but they can be renewed
func (algoInteractor *AlgoInteractor) getWalletHandleToken(walletID string, walletPassword string) (walletHandleToken string, error error) {
	initResponse, err := algoInteractor.kmdClient.InitWalletHandle(walletID, walletPassword)
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return "", err
	}
	return initResponse.WalletHandleToken, nil

}

func (algoInteractor *AlgoInteractor) getTransactions(in *pb.GetTransactionsRequest) (models.TransactionsResponse, error) {
	return algoInteractor.indexerClient.LookupAccountTransactions(in.Address).Do(context.Background())
}
