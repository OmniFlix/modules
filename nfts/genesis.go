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
		
		for _, dnft := range genState.DNFTs {
			k.SetDNFT(ctx, dnft)
		}
		
		for _, stream := range genState.LiveStreams {
			k.SetLiveStream(ctx, stream)
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
		
		for _, adnft := range genState.AdNFTs {
			count := k.GetGlobalAdsCount(ctx)
			k.MintAdNFT(ctx, adnft)
			addr, _ := sdk.AccAddressFromBech32(adnft.Owner)
			k.SetAdNFTIDToAccount(ctx, addr, adnft.AssetID)
			k.SetGlobalAdsCount(ctx, count+1)
		}
		
		for _, info := range genState.TwitterAccountInfo {
			k.SetTwitterHandlerInfo(ctx, info)
		}
		
		for _, addr := range genState.ACLAddressList.AccessList {
			k.UpdateAclAddress(ctx, addr)
		}
		
		k.UpdaterHandlerInfo(ctx, genState.ACLHandlersInfo)
		
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	
	nfts := k.GetAllTweetNFTs(ctx)
	adnfts := k.GetAllAdNFTs(ctx)
	dnfts := k.GetAllDnfts(ctx)
	liveStreams := k.GetLiveStreams(ctx)
	twitterAccountInfos := k.GetAllTwitterHandleInfos(ctx)
	aclAdresses := k.GetAclAddressList(ctx)
	aclHandlers := k.GetAuthorisedHandlerInfo(ctx)
	
	return GenesisState{
		TweetNFTs:          nfts,
		AdNFTs:             adnfts,
		DNFTs:              dnfts,
		LiveStreams:        liveStreams,
		TwitterAccountInfo: twitterAccountInfos,
		ACLHandlersInfo:    aclHandlers,
		ACLAddressList:     aclAdresses,
	}
}
