package keeper

import (
	"fmt"
	"sort"
	
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	bankKeeper    types.BaseBankKeeper
	accountKeeper types.AccountKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey,
	accountKeeper types.AccountKeeper, bankKeeper types.BaseBankKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (keeper Keeper) GetGlobalAdsCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	
	key := types.GetGlobalAdsCountKey()
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	
	var count uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	return count
}
func (keeper Keeper) SetGlobalAdsCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(keeper.storeKey)
	key := types.GetGlobalAdsCountKey()
	
	store.Set(key, keeper.cdc.MustMarshalBinaryLengthPrefixed(count))
}

func (keeper Keeper) MintAdNFT(ctx sdk.Context, nft types.BaseAdNFT) {
	key := types.GetAdsNFTKey([]byte(nft.AdNFTID))
	store := ctx.KVStore(keeper.storeKey)
	store.Set(key, keeper.cdc.MustMarshalBinaryLengthPrefixed(nft))
}

func (keeper Keeper) GetAdNFTByID(ctx sdk.Context, id string) (types.BaseAdNFT, bool) {
	store := ctx.KVStore(keeper.storeKey)
	
	bz := store.Get(types.GetAdsNFTKey([]byte(id)))
	if bz == nil {
		return types.BaseAdNFT{}, false
	}
	
	var nft types.BaseAdNFT
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &nft)
	
	return nft, true
}

func (keeper Keeper) SetAdNFTIDToAccount(ctx sdk.Context, addr sdk.AccAddress, id string) {
	adIDs := keeper.GetAdNFTIDsOfAccount(ctx, addr)
	adIDs = append(adIDs, id)
	
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.GetAdsCountOfAddressKey(addr), keeper.cdc.MustMarshalBinaryLengthPrefixed(adIDs))
}

func (keeper Keeper) GetAdNFTIDsOfAccount(ctx sdk.Context, addr sdk.AccAddress) []string {
	store := ctx.KVStore(keeper.storeKey)
	
	key := types.GetAdsCountOfAddressKey(addr.Bytes())
	bz := store.Get(key)
	if bz == nil {
		return []string{}
	}
	
	var adIDs []string
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &adIDs)
	return adIDs
}

func (keeper Keeper) GetAdsOfAccount(ctx sdk.Context, address sdk.AccAddress) []types.BaseAdNFT {
	var nfts []types.BaseAdNFT
	adIDs := keeper.GetAdNFTIDsOfAccount(ctx, address)
	
	for _, ad := range adIDs {
		var nft types.BaseAdNFT
		nft, _ = keeper.GetAdNFTByID(ctx, ad)
		nfts = append(nfts, nft)
	}
	
	return nfts
}

func (keeper Keeper) GetGlobalTweetCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	
	key := types.GetGlobalTweetCountKey()
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	
	var count uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	return count
}
func (keeper Keeper) SetGlobalTweetCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(keeper.storeKey)
	key := types.GetGlobalTweetCountKey()
	
	store.Set(key, keeper.cdc.MustMarshalBinaryLengthPrefixed(count))
}

func (k Keeper) UpdateTweetNFT(ctx sdk.Context, nft types.BaseTweetNFT) {
	k.MintTweetNFT(ctx, nft)
}

func (keeper Keeper) MintTweetNFT(ctx sdk.Context, nft types.BaseTweetNFT) {
	var key []byte
	if types.GetContextOfCurrentChain() == types.FreeFlixContext {
		key = types.GetTweetNFTKey([]byte(nft.PrimaryNFTID))
	} else {
		key = types.GetTweetNFTKey([]byte(nft.SecondaryNFTID))
	}
	
	store := ctx.KVStore(keeper.storeKey)
	store.Set(key, keeper.cdc.MustMarshalBinaryLengthPrefixed(nft))
}

func (keeper Keeper) GetTweetNFTByID(ctx sdk.Context, id string) (types.BaseTweetNFT, bool) {
	store := ctx.KVStore(keeper.storeKey)
	
	bz := store.Get(types.GetTweetNFTKey([]byte(id)))
	if bz == nil {
		return types.BaseTweetNFT{}, false
	}
	
	var nft types.BaseTweetNFT
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &nft)
	
	return nft, true
}

func (keeper Keeper) SetTweetIDToAccount(ctx sdk.Context, addr sdk.AccAddress, id string) {
	tweetIDs := keeper.GetTweetIDsOfAccount(ctx, addr)
	tweetIDs = append(tweetIDs, id)
	
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.GetTweetsCountOfAddressKey(addr), keeper.cdc.MustMarshalBinaryLengthPrefixed(tweetIDs))
}

func (keeper Keeper) GetTweetIDsOfAccount(ctx sdk.Context, addr sdk.AccAddress) []string {
	store := ctx.KVStore(keeper.storeKey)
	
	key := types.GetTweetsCountOfAddressKey(addr.Bytes())
	bz := store.Get(key)
	if bz == nil {
		return []string{}
	}
	
	var tweetIDs []string
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &tweetIDs)
	return tweetIDs
}

func (keeper Keeper) GetTweetsOfAccount(ctx sdk.Context, address sdk.AccAddress) []types.BaseTweetNFT {
	var nfts []types.BaseTweetNFT
	tweetIDs := keeper.GetTweetIDsOfAccount(ctx, address)
	
	for _, tweet := range tweetIDs {
		var nft types.BaseTweetNFT
		nft, _ = keeper.GetTweetNFTByID(ctx, tweet)
		nfts = append(nfts, nft)
	}
	
	return nfts
}

func (keeper Keeper) GetAllTweetNFTs(ctx sdk.Context) []types.BaseTweetNFT {
	store := ctx.KVStore(keeper.storeKey)
	
	iterator := sdk.KVStorePrefixIterator(store, types.TweetNFTPrefix)
	defer iterator.Close()
	
	var nfts []types.BaseTweetNFT
	for ; iterator.Valid(); iterator.Next() {
		var nft types.BaseTweetNFT
		value := iterator.Value()
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &nft)
		nfts = append(nfts, nft)
	}
	
	return nfts
}

func (keeper Keeper) SetDNFT(ctx sdk.Context, dnft types.BaseDNFT) {
	store := ctx.KVStore(keeper.storeKey)
	
	slotTime, programmeTime := types.GetTimeSlotFromDNFTID(dnft.DNFTID)
	storeKey := types.GetDnftStoreKey([]byte(slotTime), []byte(programmeTime))
	
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(dnft)
	store.Set(storeKey, bz)
}

func (keeper Keeper) GetDNFT(ctx sdk.Context, id string) (types.BaseDNFT, bool) {
	store := ctx.KVStore(keeper.storeKey)
	slotTime, programmeTime := types.GetTimeSlotFromDNFTID(id)
	
	storeKey := types.GetDnftStoreKey([]byte(slotTime), []byte(programmeTime))
	
	bz := store.Get(storeKey)
	if bz == nil {
		return types.BaseDNFT{}, false
	}
	
	var dnft types.BaseDNFT
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dnft)
	return dnft, true
}

func (keeper Keeper) SetLiveStream(ctx sdk.Context, liveStream types.BaseLiveStream) {
	store := ctx.KVStore(keeper.storeKey)
	
	storeKey := types.GetLiveStreamKey([]byte(liveStream.LiveStreamID))
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(liveStream)
	store.Set(storeKey, bz)
}

func (keeper Keeper) GetLiveStream(ctx sdk.Context, id string) (types.BaseLiveStream, bool) {
	store := ctx.KVStore(keeper.storeKey)
	storeKey := types.GetLiveStreamKey([]byte(id))
	bz := store.Get(storeKey)
	if bz == nil {
		return types.BaseLiveStream{}, false
	}
	
	var liveStream types.BaseLiveStream
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &liveStream)
	return liveStream, true
}

func (keeper Keeper) SetGlobalLiveStreamCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(keeper.storeKey)
	storekey := types.GetGlobalLiveStreamCountKey()
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(count)
	store.Set(storekey, bz)
}

func (keeper Keeper) GetGlobalLiveStreamCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	storeKey := types.GetGlobalLiveStreamCountKey()
	bz := store.Get(storeKey)
	if bz == nil {
		return 0
	}
	
	var count uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	return count
}

func (keeper Keeper) GetLiveStreams(ctx sdk.Context) []types.BaseLiveStream {
	store := ctx.KVStore(keeper.storeKey)
	
	iterator := sdk.KVStorePrefixIterator(store, types.LiveStreamPrefix)
	iterator.Close()
	
	var liveStreams []types.BaseLiveStream
	for ; iterator.Valid(); iterator.Next() {
		var liveStream types.BaseLiveStream
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &liveStream)
		
		sort.Strings(liveStream.DNFTIDs)
		liveStreams = append(liveStreams, liveStream)
	}
	return liveStreams
}

func (keeper Keeper) GetDNFTsFromIterator(ctx sdk.Context, startTime, endTime string) []types.BaseDNFT {
	store := ctx.KVStore(keeper.storeKey)
	dnftIterator := store.Iterator(types.GetDNFTTimeKey([]byte(startTime)), types.GetDNFTTimeKey([]byte(endTime)))
	defer dnftIterator.Close()
	
	var dnfts []types.BaseDNFT
	for ; dnftIterator.Valid(); dnftIterator.Next() {
		value := dnftIterator.Value()
		var dnft types.BaseDNFT
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &dnft)
		dnfts = append(dnfts, dnft)
	}
	return dnfts
}

func (k Keeper) SetTwitterHandlerInfo(ctx sdk.Context, info types.TwitterAccountInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHandlerKey([]byte(info.Handle))
	
	store.Set(key, k.cdc.MustMarshalBinaryLengthPrefixed(info))
}

func (k Keeper) GetTwitterHandleInfo(ctx sdk.Context, handle string) (types.TwitterAccountInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	
	key := types.GetHandlerKey([]byte(handle))
	bz := store.Get(key)
	if bz == nil {
		return types.TwitterAccountInfo{}, false
	}
	
	var info types.TwitterAccountInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	return info, true
}

func (k Keeper) GetAllTwitterHandleInfos(ctx sdk.Context) []types.TwitterAccountInfo {
	store := ctx.KVStore(k.storeKey)
	
	iterator := sdk.KVStorePrefixIterator(store, types.TwitterHandlePrefix)
	defer iterator.Close()
	
	var infos []types.TwitterAccountInfo
	for ; iterator.Valid(); iterator.Next() {
		var info types.TwitterAccountInfo
		value := iterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &info)
		infos = append(infos, info)
	}
	
	return infos
}

func (k Keeper) SendCoins(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from, to, amount)
}

func (keeper Keeper) GetAllAdNFTs(ctx sdk.Context) []types.BaseAdNFT {
	store := ctx.KVStore(keeper.storeKey)
	
	iterator := sdk.KVStorePrefixIterator(store, types.AdNFTPrefix)
	defer iterator.Close()
	
	var adnfts []types.BaseAdNFT
	for ; iterator.Valid(); iterator.Next() {
		var adnft types.BaseAdNFT
		
		value := iterator.Value()
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &adnft)
		adnfts = append(adnfts, adnft)
	}
	
	return adnfts
}

func (keeper Keeper) GetAllDnfts(ctx sdk.Context) []types.BaseDNFT {
	store := ctx.KVStore(keeper.storeKey)
	
	iterator := sdk.KVStorePrefixIterator(store, types.DNFTPrefix)
	defer iterator.Close()
	
	var dnfts []types.BaseDNFT
	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		
		var dnft types.BaseDNFT
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &dnft)
		dnfts = append(dnfts, dnft)
	}
	
	return dnfts
}

func (k Keeper) UpdateAclAddress(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	addrs := k.GetAclAddressList(ctx)
	addrs.AccessList = append(addrs.AccessList, addr)
	
	store.Set(types.GetAclKey(), k.cdc.MustMarshalBinaryLengthPrefixed(addrs))
}

func (k Keeper) GetAclAddressList(ctx sdk.Context) types.AclInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAclKey())
	if bz == nil {
		return types.AclInfo{}
	}
	
	var addresses types.AclInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &addresses)
	return addresses
}

func (k Keeper) UpdaterHandlerInfo(ctx sdk.Context, handlers types.AllowedHandles) {
	store := ctx.KVStore(k.storeKey)
	
	store.Set(types.GetAuthorisedHandlersKey(), k.cdc.MustMarshalBinaryLengthPrefixed(handlers))
}

func (k Keeper) GetAuthorisedHandlerInfo(ctx sdk.Context) types.AllowedHandles {
	store := ctx.KVStore(k.storeKey)
	
	bz := store.Get(types.GetAuthorisedHandlersKey())
	if bz == nil {
		return types.AllowedHandles{}
	}
	var handles types.AllowedHandles
	
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &handles)
	handles.Sort()
	return handles
}

func (k Keeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	return k.accountKeeper.GetAccount(ctx, addr)
}
