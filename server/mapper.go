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
