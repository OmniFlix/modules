package keeper

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abcitypes.RequestQuery) (bytes []byte, err error) {
		switch path[0] {
		case types.QueryTweetNFT:
			return queryUsingNFTID(ctx, path[1:], k)
		case types.QueryTweetNFTsByAddress:
			return queryTweetNFTsByAddress(ctx, path[1:], k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
		
	}
}

func queryUsingNFTID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	nft, found := k.GetTweetNFTByID(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, fmt.Sprintf("nft %s ", path[0]))
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, nft)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}

func queryTweetNFTsByAddress(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("nft %s ", path[0]))
	}
	
	tweeets := k.GetTweetsOfAccount(ctx, addr)
	
	res, err := codec.MarshalJSONIndent(k.cdc, tweeets)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}
