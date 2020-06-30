package keeper

import (
	"reflect"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/abci/types"
	
	types2 "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req types.RequestQuery) (bytes []byte, err error) {
		switch path[0] {
		case types2.QueryParams:
			return queryParams(ctx, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
		
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params, _ := k.GetParams(ctx)
	
	if reflect.DeepEqual(params, types2.Params{}) {
		return nil, sdkerrors.Wrap(types2.ErrParamsNotFound, "")
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}
