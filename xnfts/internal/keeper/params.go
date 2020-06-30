package keeper

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params, found bool) {
	store := ctx.KVStore(k.storeKey)
	
	bz := store.Get(types.GetParamKey())
	
	if bz == nil {
		return types.Params{}, false
	}
	
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &params)
	return params, true
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	
	store.Set(types.GetParamKey(), k.cdc.MustMarshalBinaryLengthPrefixed(params))
}

func (k Keeper) GetLastVistedTime(ctx sdk.Context) time.Time {
	store := ctx.KVStore(k.storeKey)
	
	bz := store.Get(types.GetLastVisitedKey())
	if bz == nil {
		return time.Time{}
	}
	
	var timeData time.Time
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &timeData)
	return timeData
}

func (k Keeper) SetLastVisitedTime(ctx sdk.Context, timeBytes []byte) {
	store := ctx.KVStore(k.storeKey)
	
	store.Set(types.GetLastVisitedKey(), k.cdc.MustMarshalBinaryLengthPrefixed(timeBytes))
}
