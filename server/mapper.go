package main

import (
	pb "github.com/LorenzHW/algo-grpc/protos"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func mapToAccount(algoAccount models.Account) *pb.Account {
	account := pb.Account{
		Amount:         algoAccount.Amount,
		Address:        algoAccount.Address,
		Assets:         mapToAssets(algoAccount.Assets),
		Rewards:        algoAccount.Rewards,
		CreatedAtRound: algoAccount.CreatedAtRound,
	}
	return &account
}

func mapToAssets(algoAssets []models.AssetHolding) map[int32]*pb.AssetHolding {
	assets := make(map[int32]*pb.AssetHolding)
	for key, value := range algoAssets {
		assetHolding := pb.AssetHolding{
			Creator: value.Creator,
			Amount:  value.Amount,
			Frozen:  value.IsFrozen,
		}
		assets[int32(key)] = &assetHolding
	}
	return assets
}

func mapToTransactions(algodTransactions []models.Transaction) []*pb.Transaction {
	transactions := make([]*pb.Transaction, 0)

	for _, algoTransaction := range algodTransactions {
		t := pb.Transaction{
			From:   algoTransaction.Sender,
			To:     algoTransaction.PaymentTransaction.Receiver,
			Amount: algoTransaction.PaymentTransaction.Amount,
			Fee:    algoTransaction.Fee,
		}
		transactions = append(transactions, &t)
	}
	return transactions
}

func mapToBlock(algoBlock models.Block) *pb.Block {
	block := pb.Block{
		Round:        algoBlock.Round,
		Transactions: mapToTransactions(algoBlock.Transactions),
		Timestamp:    0,
	}
	return &block
}
