package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

func (k Keeper) GetTweetNFTByID(ctx sdk.Context, id string) (nfts.BaseTweetNFT, bool) {
	return k.nftKeeper.GetTweetNFTByID(ctx, id)
}

func (k Keeper) GetAllTweetNFTs(ctx sdk.Context) []nfts.BaseTweetNFT {
	return k.nftKeeper.GetAllTweetNFTs(ctx)
}

func (k Keeper) GetAdNFTByID(ctx sdk.Context, id string) (nfts.BaseAdNFT, bool) {
	return k.nftKeeper.GetAdNFTByID(ctx, id)
}

func (k Keeper) MintAdNFTByID(ctx sdk.Context, nft nfts.BaseAdNFT) {
	k.nftKeeper.MintAdNFT(ctx, nft)
}

func (k Keeper) GetGlobalTweetCount(ctx sdk.Context) uint64 {
	return k.nftKeeper.GetGlobalTweetCount(ctx)
}

func (k Keeper) SetGlobalTweetCount(ctx sdk.Context, count uint64) {
	k.nftKeeper.SetGlobalTweetCount(ctx, count)
}

func (k Keeper) GetTweetsOfAccount(ctx sdk.Context, address sdk.AccAddress) []nfts.BaseTweetNFT {
	return k.nftKeeper.GetTweetsOfAccount(ctx, address)
}

func (k Keeper) SetTweetIDToAccount(ctx sdk.Context, addr sdk.AccAddress, id string) {
	k.nftKeeper.SetTweetIDToAccount(ctx, addr, id)
}

func (k Keeper) MintTweetNFT(ctx sdk.Context, nft nfts.BaseTweetNFT) {
	k.nftKeeper.MintTweetNFT(ctx, nft)
}

func (k Keeper) GetTwitterHandleInfo(ctx sdk.Context, handle string) (nfts.TwitterAccountInfo, bool) {
	return k.nftKeeper.GetTwitterHandleInfo(ctx, handle)
}

func (k Keeper) SetTwitterHandlerInfo(ctx sdk.Context, info nfts.TwitterAccountInfo) {
	k.nftKeeper.SetTwitterHandlerInfo(ctx, info)
}

func (k Keeper) GetDNFT(ctx sdk.Context, programTime string) {
	k.nftKeeper.GetDNFT(ctx, programTime)
}

func (k Keeper) SetDNFT(ctx sdk.Context, dnft nfts.BaseDNFT) {
	k.nftKeeper.SetDNFT(ctx, dnft)
}

func (k Keeper) GetLiveStream(ctx sdk.Context, id string) {
	k.nftKeeper.GetLiveStream(ctx, id)
}

func (k Keeper) SetLiveStream(ctx sdk.Context, liveStream nfts.BaseLiveStream) {
	k.nftKeeper.SetLiveStream(ctx, liveStream)
}

func (k Keeper) GetAllLiveStreams(ctx sdk.Context) []nfts.BaseLiveStream {
	return k.nftKeeper.GetLiveStreams(ctx)
}

func (k Keeper) GetDNFTsBetweenInterval(ctx sdk.Context, startTime string, endTime string) []nfts.BaseDNFT {
	return k.nftKeeper.GetDNFTsFromIterator(ctx, startTime, endTime)
}

func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error) {
	return k.bankKeeper.AddCoins(ctx, addr, amount)
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error) {
	return k.bankKeeper.SubtractCoins(ctx, addr, amount)
}

func (k Keeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	return k.accountKeeper.GetAccount(ctx, addr)
}
