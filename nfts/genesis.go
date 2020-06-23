package nfts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, genState GenesisState) {
	if GetContextOfCurrentChain() == CoCoContext {
		for _, nft := range genState.TweetNFTs {
			count := k.GetGlobalTweetCount(ctx)
			addr, _ := sdk.AccAddressFromBech32(nft.SecondaryOwner)
			k.MintTweetNFT(ctx, nft)
			k.SetTweetIDToAccount(ctx, addr, nft.SecondaryNFTID)
			k.SetGlobalTweetCount(ctx, count+1)
		}
	}
	
	if GetContextOfCurrentChain() == FreeFlixContext {
		for _, nft := range genState.TweetNFTs {
			count := k.GetGlobalTweetCount(ctx)
			addr, _ := sdk.AccAddressFromBech32(nft.PrimaryOwner)
			k.MintTweetNFT(ctx, nft)
			k.SetTweetIDToAccount(ctx, addr, nft.PrimaryNFTID)
			k.SetGlobalTweetCount(ctx, count+1)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	nfts := k.GetAllTweetNFTs(ctx)
	
	return GenesisState{
		TweetNFTs: nfts,
	}
}
