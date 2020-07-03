

### Querier

To query the data of our app, we need to make it accessible using our `Querier`. This piece of the app works in tandem with the `Keeper` to access the state and return it. The `Querier` is defined in `./nfts/internal/keeper/querier.go`. Our `scaffold` tool starts us out with some suggestions on how it should look, and similar to our `Handler` we want to handle different queried routes. You could make many different routes within the `Querier` for many different types of queries, but we will just make three:
*  `queryUsingNFTID`- query tweetNFT using nftID
*  `queryTweetNFTsByAddress`- query all tweetNFTs for an address

Combined into a switch statement and with each of the functions fleshed out it should look as follows:

```go=
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

```



**Types**

You may notice that we use two different imported types on our initial switch statement. These are defined within our `./nfts/internal/types/querier.go ` file as simple strings. That file should look like the following:

```go=
package types

const (
    QueryTweetNFT           = "tweet_nft"
    QueryTweetNFTsByAddress = "address_tweet_nfts"
)

```