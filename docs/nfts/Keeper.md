
### Keeper

After using the `scaffold` command you should have a boilerplate `Keeper` at ./nfts/internal/keeper/keeper.go. It contains a basic keeper with references to basic functions like `Set`, `Get` and `Delete`.

Our keeper stores all our data for our module. Sometimes a module will import the keeper of another module. This will allow the state to be shared and modified across modules. Here we are dealing with the creation of nfts. Look at our completed `Keeper` and how `Set` and `Get` are expanded:

```go=
package keeper

import (
    "fmt"
    
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/tendermint/tendermint/libs/log"
    
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

type Keeper struct {
    storeKey sdk.StoreKey
    cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
    return Keeper{
        storeKey: key,
        cdc:      cdc,
    }
}

func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
    return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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

```