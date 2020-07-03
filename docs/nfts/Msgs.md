
### Messages

Messages are a great place to start when building a module because they define the actions that your application can make. Think of all the scenarios where a user would be able to update the state of the application in any way. These should be boiled down into basic interactions, similar to **CRUD** (Create, Read, Update and Delete).

Let's start with **Create**

**MsgMintTweetNFT**

Messages are `types` which live inside the `./nfts/internal/types/` directory. There is already a `msg.go` file. We can use this`msg.go`.

```go=
package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgMintTweetNFT struct {
    Sender        sdk.AccAddress `json:"sender"`
    AssetID       string         `json:"asset_id"`
    License       bool           `json:"license"`
    LicensingFee  sdk.Coin       `json:"licensing_fee"`
    RevenueShare  sdk.Dec        `json:"revenue_share"`
    TwitterHandle string         `json:"twitter_handle"`
}

func NewMsgMintNFT(sender sdk.AccAddress, assetID string, license bool, fee sdk.Coin, share sdk.Dec, handle string) MsgMintTweetNFT {
    return MsgMintTweetNFT{
        Sender:        sender,
        AssetID:       assetID,
        License:       license,
        LicensingFee:  fee,
        RevenueShare:  share,
        TwitterHandle: handle,
    }
}

var _ sdk.Msg = MsgMintTweetNFT{}

func (m MsgMintTweetNFT) Route() string {
    return RouterKey
}

func (m MsgMintTweetNFT) Type() string {
    return "msg_mint_tweet_nft"
}

func (m MsgMintTweetNFT) ValidateBasic() error {
    
    if m.Sender.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
    } else if m.AssetID == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
    }
    
    if m.License {
        if m.LicensingFee.IsZero() {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid licensing fee provided")
        } else if m.RevenueShare.IsZero() {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share should not be nil")
        }
    }
    if m.TwitterHandle == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "twitter handle should not be empty")
    }
    return nil
}

func (m MsgMintTweetNFT) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgMintTweetNFT) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{m.Sender}
}

```

Notice that all Messages in the app need to follow the `sdk.Msg` interface. The Message `struct` contains all the necessary information when creating a new NFT:

* `Sender`- The account that creates a new NFT. This uses the `sdk.AccAddress` type which represents an account in the app. 
* `AssetID`- AssetID of the `Tweet NFT`
* `Licence`- Tweet can be licensed by another account, if this is set to `TRUE`. It can be Un-licensed if it is set to `FALSE`.
* `LicensingFee`- Amount to pay, to license an asset
* `RevenueShare`- Amount to share with the owner when a licensee distributes (or utilizes) an asset.
* `TwitterHandle`- Twitter handle of the user

The `Msg` interface requires some other methods to be set, like validating the content of the `struct`, and confirming the msg was signed and submitted by the Sender.