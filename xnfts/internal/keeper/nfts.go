package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

func (k Keeper) GetTweetNFTByID(ctx sdk.Context, id string) (nfts.BaseTweetNFT, bool) {
	return k.nftKeeper.GetTweetNFTByID(ctx, id)
}

func (k Keeper) GetAllTweetNFTs(ctx sdk.Context) []nfts.BaseTweetNFT {
	return k.nftKeeper.GetAllTweetNFTs(ctx)
}

func (k Keeper) GetGlobalTweetCount(ctx sdk.Context) uint64 {
	return k.nftKeeper.GetGlobalTweetCount(ctx)
}

func (k Keeper) SetGlobalTweetCount(ctx sdk.Context, count uint64) {
	k.nftKeeper.SetGlobalTweetCount(ctx, count)
	return
}

func (k Keeper) MintTweetNFT(ctx sdk.Context, nft nfts.BaseTweetNFT) {
	k.nftKeeper.MintTweetNFT(ctx, nft)
	return
}

func (k Keeper) SetTweetIDToAccount(ctx sdk.Context, addr sdk.AccAddress, id string) {
	k.nftKeeper.SetTweetIDToAccount(ctx, addr, id)
	return
}
