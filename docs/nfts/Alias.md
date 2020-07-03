
### Alias
Now that we have these new message types, we'd like to make sure other parts of the module can access them. To do so we use the `./nfts/alias.go` file. This imports the types from the nested types directory and makes them accessible at the modules top level directory.

```go=
package nfts

import (
    "github.com/FreeFlixMedia/modules/nfts/internal/keeper"
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

const (
    CoCoContext     = types.CoCoContext
    FreeFlixContext = types.FreeFlixContext
    
    ModuleName = types.ModuleName
    RouterKey  = types.RouterKey
    StoreKey   = types.StoreKey
)

type (
    Keeper       = keeper.Keeper
    GenesisState = types.GenesisState
    
    MsgMintTweetNFT = types.MsgMintTweetNFT
    BaseTweetNFT    = types.BaseTweetNFT
)

var (
    NewKeeper                = keeper.NewKeeper
    NewQuerier               = keeper.NewQuerier
    GetContextOfCurrentChain = types.GetContextOfCurrentChain
    GetPrimaryNFTID          = types.GetPrimaryNFTID
    GetSecondaryNFTID        = types.GetSecondaryNFTID
    
    EventTypeMsgMintTweetNFT = types.EventTypeMsgMintTweetNFT
    
    AttributePrimaryNFTID   = types.AttributePrimaryNFTID
    AttributeSecondaryNFTID = types.AttributeSecondaryNFTID
    AttributeAssetID        = types.AttributeAssetID
    AttributeTwitterHandle  = types.AttributeTwitterHandle
    
    ErrAssetIDAlreadyExist = types.ErrAssetIDAlreadyExist
    ErrInvalidLicense      = types.ErrInvalidLicense
    ErrParamsNotFound      = types.ErrParamsNotFound
    ErrNFTNotFound         = types.ErrNFTNotFound
)

```

It's great to have Messages, but we need somewhere to store the information they are sending. All persistent data related to this module should live in the module's `Keeper`.

Let's make a `Keeper` for our NFTs Module next.