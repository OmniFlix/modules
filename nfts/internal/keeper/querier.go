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
		case types.QueryLiveStream:
			return queryLiveStream(ctx, path[1:], k)
		case types.QueryDNFT:
			return queryDNFT(ctx, path[1:], k)
		case types.QueryLiveStreams:
			return queryLiveStreams(ctx, k)
		case types.QueryAdNFTBYID:
			return queryAdNFTByID(ctx, path[1:], k)
		case types.QueryAdNFTByAddress:
			return queryAdNFTsByAddress(ctx, path[1:], k)
		case types.QueryTwitterAccount:
			return queryTwitterAccountInfoByHandle(ctx, path[1:], k)
		case types.QueryTwitterAccounts:
			return queryAllTwitterAccountsInfo(ctx, k)
		case types.QueryAuthorisedHandler:
			return queryAuthorisedHandlersInfo(ctx, k)
		case types.QueryAuthorisedAddresses:
			return queryAuthorisedAddresses(ctx, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
		
	}
}

func queryUsingNFTID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	nft, found := k.GetTweetNFTByID(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, fmt.Sprintf("nfts %s ", path[0]))
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("nfts %s ", path[0]))
	}
	
	tweeets := k.GetTweetsOfAccount(ctx, addr)
	
	res, err := codec.MarshalJSONIndent(k.cdc, tweeets)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}

func queryDNFT(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	
	liveStream, found := keeper.GetDNFT(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, fmt.Sprintf("dnft %s ", path[0]))
	}
	
	bz, err := codec.MarshalJSONIndent(keeper.cdc, liveStream)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return bz, nil
}

func queryLiveStream(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	
	liveStream, found := keeper.GetLiveStream(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, fmt.Sprintf("liveStream %s ", path[0]))
	}
	
	bz, err := codec.MarshalJSONIndent(keeper.cdc, liveStream)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return bz, nil
}

func queryLiveStreams(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	streams := keeper.GetLiveStreams(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, streams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return bz, nil
}

func queryAdNFTByID(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	adNFT, found := keeper.GetAdNFTByID(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, fmt.Sprintf("adnft %s ", path[0]))
	}
	
	bz, err := codec.MarshalJSONIndent(keeper.cdc, adNFT)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return bz, nil
	
}

func queryAdNFTsByAddress(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("nfts %s ", path[0]))
	}
	
	ads := k.GetAdsOfAccount(ctx, addr)
	
	res, err := codec.MarshalJSONIndent(k.cdc, ads)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}

func queryTwitterAccountInfoByHandle(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	info, found := k.GetTwitterHandleInfo(ctx, path[0])
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFound, fmt.Sprintf("handle %s ", path[0]))
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, info)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}

func queryAllTwitterAccountsInfo(ctx sdk.Context, k Keeper) ([]byte, error) {
	accounts := k.GetAllTwitterHandleInfos(ctx)
	res, err := codec.MarshalJSONIndent(k.cdc, accounts)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
}

func queryAuthorisedAddresses(ctx sdk.Context, k Keeper) ([]byte, error) {
	address := k.GetAclAddressList(ctx)
	
	res, err := codec.MarshalJSONIndent(k.cdc, address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
	
}

func queryAuthorisedHandlersInfo(ctx sdk.Context, k Keeper) ([]byte, error) {
	address := k.GetAuthorisedHandlerInfo(ctx)
	
	res, err := codec.MarshalJSONIndent(k.cdc, address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	
	return res, nil
	
}
