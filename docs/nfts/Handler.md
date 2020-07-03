

### Handler

For a **Message** to reach a **Keeper**, it has to go through a **Handler**. This is where logic can be applied to either allow or deny a `Message` to succeed. It's also where logic as to exactly how the state should change within the Keeper should take place. If you're familiar with [Model View Controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller) (MVC) architecture, the `Keeper` is a bit like the **Model** and the `Handler` is a bit like the **Controller**. If you're familiar with [React/Redux](https://en.wikipedia.org/wiki/React_(web_framework)) or [Vue/Vuex](https://en.wikipedia.org/wiki/Vue.js) architecture, the `Keeper` is a bit like the **Reducer/Store** and the `Handler` is a bit like **Actions**.

Our Handler will go in `./nfts/handler.go` and will follow the suggestions outlined in the boilerplate. We will create handler functions for `Message` type, `MsgMintTweetNFT` until the file looks as follows:

```go=
package nfts

import (
    "strings"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
        ctx = ctx.WithEventManager(sdk.NewEventManager())
        
        switch msg := msg.(type) {
        case MsgMintTweetNFT:
            return handleMsgMintTweetNFT(ctx, keeper, msg)
        
        default:
            return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized NFT message type: %T", msg)
        }
    }
}

func handleMsgMintTweetNFT(ctx sdk.Context, keeper Keeper, msg MsgMintTweetNFT) (*sdk.Result, error) {
    
    nfts := keeper.GetTweetsOfAccount(ctx, msg.Sender)
    
    for _, nft := range nfts {
        if strings.EqualFold(nft.AssetID, msg.AssetID) {
            return nil, sdkerrors.Wrap(ErrAssetIDAlreadyExist, "")
        }
    }
    
    count := keeper.GetGlobalTweetCount(ctx)
    id := GetPrimaryNFTID(count)
    tweetNFT := BaseTweetNFT{
        PrimaryNFTID:   id,
        PrimaryOwner:   msg.Sender.String(),
        SecondaryNFTID: "",
        SecondaryOwner: "",
        License:        msg.License,
        AssetID:        msg.AssetID,
        LicensingFee:   msg.LicensingFee,
        RevenueShare:   msg.RevenueShare,
        TwitterHandle:  msg.TwitterHandle,
    }
    
    keeper.MintTweetNFT(ctx, tweetNFT)
    keeper.SetTweetIDToAccount(ctx, msg.Sender, tweetNFT.PrimaryNFTID)
    keeper.SetGlobalTweetCount(ctx, count+1)
    
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            EventTypeMsgMintTweetNFT,
            sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
            sdk.NewAttribute(AttributePrimaryNFTID, tweetNFT.PrimaryNFTID),
            sdk.NewAttribute(AttributeAssetID, tweetNFT.AssetID),
            sdk.NewAttribute(AttributeTwitterHandle, tweetNFT.TwitterHandle),
        ),
    )
    
    return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
    
}

```

#### GetGlobalTweetCount & SetGlobalTweetCount

We use `SetGlobalTweetCount` to increase count value when new TweetNFT comes. And we use `GetGlobalTweetCount` to get a count of all tweets that are created.