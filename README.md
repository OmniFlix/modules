Modules
===

I) NFT Module
---

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

### Codec

Once we have defined our messages, we need to describe to our encoder how they should be stored as bytes. To do this we edit the file located at `./nfts/internals/types/codec.go`. By describing our types as follows they will work with our encoding library:

```go=
package types

import (
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.Codec) {
    cdc.RegisterConcrete(MsgMintTweetNFT{}, "nft/MsgMintTweetNFT", nil)
    cdc.RegisterConcrete(BaseTweetNFT{}, "nft/BaseTweetNFT", nil)
}

var (
    amino     = codec.New()
    ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
    RegisterCodec(amino)
    amino.Seal()
}

```

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

### BaseTweetNFT

You may notice the reference to `types.BaseTweetNFT` throughout the `Keeper`. These is new struct defined in `./nfts/types/nfts.go` that contains all necessary information about different tweet nfts. You can create this file now and add the following:

```go=
package types

import (
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type BaseTweetNFT struct {
    PrimaryNFTID string `json:"primary_nft_id"`
    PrimaryOwner string `json:"primary_owner"`
    
    SecondaryNFTID string `json:"secondary_nft_id"`
    SecondaryOwner string `json:"secondary_owner"`
    
    License bool   `json:"license"`
    AssetID string `json:"asset_id"`
    
    LicensingFee sdk.Coin `json:"licensing_fee"`
    RevenueShare sdk.Dec  `json:"revenue_share"`
    
    TwitterHandle string `json:"twitter_handle"`
}

func (nft BaseTweetNFT) String() string {
    return fmt.Sprintf(`
PrimaryNFTID: %s,
PrimaryOwner: %s,

SecondaryNFTID: %s,
SecondaryOwner: %s,

License: %t,
AssetID: %s,

LicensingFee: %s,
RevenueShare: %s,

TwitterHandle: %s,
`, nft.PrimaryNFTID, nft.PrimaryOwner, nft.SecondaryNFTID, nft.SecondaryOwner,
        nft.License, nft.AssetID, nft.LicensingFee.String(), nft.RevenueShare.String(), nft.TwitterHandle)
}

```

You might also notice that each type has the `String` method. This allows us to render the struct as a string for rendering.

### Prefixes

You may notice the use of `types.FreeFlixNFTPrefix` and `types.CoCoNFTPrefix`. These are defined in a file called `./nfts/internal/types/key.go` and help us keep our `Keeper` organized. The `Keeper` is just a key-value store. That means that similar to an `Object` in javascript, all values are referenced under a key. To access a value, you need to know the key under which it is stored. This is a bit like a unique identifier (UID).

When storing a `TweetNFT` we use the key of the `PrimartNFTID` as a unique ID if it is from `FreeFlix chain`, otherwise we use the key `SecondaryNFTID` if it is from `CoCo chain`. However since we are storing it in the same location, we may want to distinguish between the types of nftIDs we use as keys. We can do this by adding prefixes to the nftIDs that allow us to recognize which is which. For `FreeFlixNFT` we add the prefix `ffmttweetnft` and for `CoCoNFT` we add the prefix `cocotweetnft`. You should add these to your `key.go` file so it looks as follows:

```go=
package types

import (
    "strconv"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
    ModuleName   = "nfts"
    RouterKey    = ModuleName
    QuerierRoute = ModuleName
    StoreKey     = ModuleName
    
    FreeFlixNFTPrefix = "ffmttweetnft"
    CoCoNFTPrefix     = "cocotweetnft"
    
    FreeFlixContext = "freeflix"
    CoCoContext     = "coco"
)

var (
    GlobalTweetCountPrefix = []byte{0x01}
    TweetAccountPrefix     = []byte{0x02}
    TweetNFTPrefix         = []byte{0x03}
)

func GetGlobalTweetCountKey() []byte {
    return GlobalTweetCountPrefix
}

func GetTweetNFTKey(id []byte) []byte {
    return append(TweetNFTPrefix, id...)
}

func GetTweetsCountOfAddressKey(addr []byte) []byte {
    return append(TweetAccountPrefix, addr...)
}

func GetPrimaryNFTID(count uint64) string {
    return FreeFlixNFTPrefix + strconv.Itoa(int(count))
}

func GetContextOfCurrentChain() string {
    config := sdk.GetConfig()
    return config.GetBech32AccountAddrPrefix()
}

func GetSecondaryNFTID(count uint64) string {
    return CoCoNFTPrefix + strconv.Itoa(int(count))
}

```

### Iterators

Sometimes you will want to access a `TweetNFT` directly by their key. That's why we have the methods `MintTweetNFT` and `GetTweetNFTByID`. However, sometimes you will want to get every `TweetNFT` at once. To do this we use an Iterator called `KVStorePrefixIterator`. This utility comes from the `sdk` and iterates over a key store. If you provide a prefix, it will only iterate over the keys that contain that prefix. Since we have prefixes defined for our `TweetNFT`, we can use them here to only return our desired data types.

---
Now that you've seen the `Keeper` where every `TweetNFT` that is stored, we need to connect the messages to this storage. This process is called *handling* the messages and is done inside the `Handler`.

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

### GetGlobalTweetCount & SetGlobalTweetCount

We use `SetGlobalTweetCount` to increase count value when new TweetNFT comes. And we use `GetGlobalTweetCount` to get a count of all tweets that are created.

### Events

At the end of each handler is an EventManager which will create logs within the transaction that reveals information about what occurred during the handling of this message. This is useful for client-side software that wants to know exactly what happened as a result of this state transition. These Events use a series of pre-defined types that can be found in `./nfts/internal/types/events.go` and look as follows:

```go=
package types

var (
    EventTypeMsgMintTweetNFT = "msg_mint_tweet_nft"
    
    AttributePrimaryNFTID   = "primary_nft_id"
    AttributeSecondaryNFTID = "secondary_nft_id"
    
    AttributeAssetID       = "asset_id"
    AttributeTwitterHandle = "twitter_handler"
)

```

Now that we have all the necessary pieces for updating state (`Message`, `Handler`, `Keeper`) we might want to consider ways in which we can query state. This is typically done via a REST endpoint and/or a CLI. Both of those clients interact with a part of the app which queries state, called the `Querier`.

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

Our queries are rather simple since we've already outfitted our `Keeper` with all the necessary functions to access the state. You can see the iterator being used here as well.

Now that we have all of the basic actions of our module created, we want to make them accessible. We can do this with a CLI client and a REST client. For this tutorial, we will just be creating a CLI client. If you are interested in what goes into making a REST client.

Let's take a look at what goes into making a CLI.

### CLI

A Command Line Interface (CLI) will help us interact with our app once it is running on a machine somewhere. Each Module has its own namespace within the CLI that gives it the ability to create and sign Messages destined to be handled by that module. It also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI will let you do things like generate keys for a new account or check the status of an interaction you already had with the application.


The CLI for our module is broken into two files called `tx.go` and `query.go` which are located in `./nfts/client/cli/.` One file is for making transactions that contain messages which will ultimately update our state. The other is for making queries which will give us the ability to read information from our state. Both files utilize the [Cobra](https://github.com/spf13/cobra) library.

**Transactions**

The `tx.go` file contains `GetTxCmd` which is a standard method within the Cosmos SDK. It is referenced later in the `module.go` file which describes exactly which attributes a module has. This makes it easier to incorporate different modules for different reasons at the level of the actual application. After all, we are focusing on a module at this point, but later we will create an application that utilizes this module as well as other modules that are already available within the Cosmos SDK.

Inside `GetTxCmd` we create a new module-specific command and call is `nfts`. Within this command we add a sub-command for each Message type we've defined:

* `GetMsgMintTweetNFT`


Each function takes parameters from the **Cobra** CLI tool to create a new msg, sign it and submit it to the application to be processed. These functions should go into the `tx.go` file and look as follows:

```go=
package cli

import (
    "bufio"
    "strconv"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
    NFTTxCmd := &cobra.Command{
        Use:   types.ModuleName,
        Short: "nfts  transfer transaction subcommands",
    }
    
    NFTTxCmd.AddCommand(flags.PostCommands(
        GetMsgMintTweetNFT(cdc),
    )...)
    
    return NFTTxCmd
}

func GetMsgMintTweetNFT(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mint-nft ",
        Short: "mint tweet nft",
        RunE: func(cmd *cobra.Command, args []string) error {
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
            cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
            
            var fee sdk.Coin
            var share sdk.Dec
            var err error
            
            license, err := strconv.ParseBool(viper.GetString(FlagLicence))
            if err != nil {
                return err
            }
            
            feeStr := viper.GetString(FlagLicenceFee)
            if feeStr != "" {
                fee, err = sdk.ParseCoin(feeStr)
                if err != nil {
                    return err
                }
                
            }
            
            shareStr := viper.GetString(FlagRevenueShare)
            if shareStr != "" {
                share, err = sdk.NewDecFromStr(shareStr)
                if err != nil {
                    return err
                }
                
            }
            
            msg := types.NewMsgMintNFT(cliCtx.GetFromAddress(), viper.GetString(FlagAssetID), license, fee, share, viper.GetString(FlagTwitterHandle))
            if err := msg.ValidateBasic(); err != nil {
                return err
            }
            return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
    
    cmd.Flags().String(FlagAssetID, "", "AssetID")
    cmd.Flags().String(FlagTwitterHandle, "", "Twitter handle")
    cmd.Flags().String(FlagLicenceFee, "0coco", "Twitter handle")
    cmd.Flags().String(FlagRevenueShare, "0", "Revenue share")
    cmd.Flags().String(FlagLicence, "false", "license")
    return cmd
}

```

**Query**

The `query.go` file contains similar **Cobra** commands that reserve a new namespace for referencing our `nfts` module. Instead of creating and submitting messages, however, the `query.go` the file creates queries and returns the results in human-readable form. The queries it handles are the same we defined in our `querier.go` file earlier:

* `GetCmdQueryTweetNFT`
* `GetCmdQueryTweetsByAccount`

After defining these commands, your `query.go` file should look like:

```go=
package cli

import (
    "fmt"
    
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/spf13/cobra"
    
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      "Querying commands for the nft module",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }
    
    cmd.AddCommand(
        GetCmdQueryTweetNFT(cdc),
        GetCmdQueryTweetsByAccount(cdc),
    )
    
    return cmd
}

func GetCmdQueryTweetNFT(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "nft [id]",
        Short: "Get NFT using nft id ",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTweetNFT, args[0]), nil)
            if err != nil {
                return err
            }
            
            var tweetNFT types.BaseTweetNFT
            cdc.MustUnmarshalJSON(res, &tweetNFT)
            return cliCtx.PrintOutput(tweetNFT)
        },
    }
    return flags.GetCommands(cmd)[0]
}

func GetCmdQueryTweetsByAccount(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "nfts [address]",
        Short: "Get  NFTs associated to account",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTweetNFTsByAddress, args[0]), nil)
            if err != nil {
                return err
            }
            
            var tweetNFTs []types.BaseTweetNFT
            cdc.MustUnmarshalJSON(res, &tweetNFTs)
            return cliCtx.PrintOutput(tweetNFTs)
        },
    }
    return flags.GetCommands(cmd)[0]
    
}

```

While these are all the major moving pieces of a module (`Message`, `Handler`, `Keeper`, `Querier`, and `Client`) there are some organizational tasks that we have yet to complete. The next step will be making sure that our module is completely configured to make it usable within any application.

### Module

Our `scaffold` tool has done most of the work for us in generating our `module.go` file inside `./nfts/.` One way that our module is different than the simplest form of a module, is that it uses it's own `Keeper`. The only real changes needed are under the `AppModule` and `NewAppModule`. The file should look as follows afterward:

```go=
package nfts

import (
    "encoding/json"
    "fmt"
    
    "github.com/gorilla/mux"
    "github.com/spf13/cobra"
    abci "github.com/tendermint/tendermint/abci/types"
    
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"
    
    "github.com/FreeFlixMedia/modules/nfts/client/cli"
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

var (
    _ module.AppModule      = AppModule{}
    _ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
    types.RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
    return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
    var data types.GenesisState
    if err := cdc.UnmarshalJSON(bz, &data); err != nil {
        return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
    }
    
    return data.ValidateGenesis()
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetQueryCmd(cdc)
}

type AppModule struct {
    AppModuleBasic
    nftKeeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
    return AppModule{
        nftKeeper: keeper,
    }
}

func (AppModule) Name() string {
    return ModuleName
}

func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() sdk.Handler { return NewHandler(am.nftKeeper) }

func (AppModule) QuerierRoute() string {
    return types.QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
    return NewQuerier(am.nftKeeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
    var genesisState GenesisState
    cdc.MustUnmarshalJSON(data, &genesisState)
    InitGenesis(ctx, am.nftKeeper, genesisState)
    return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
    gs := ExportGenesis(ctx, am.nftKeeper)
    return cdc.MustMarshalJSON(gs)
}

func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
    return []abci.ValidatorUpdate{}
}

```

Congratulations you have completed the `nfts` module!

This module is now able to be incorporated into any Cosmos SDK application.

Since we don't want to just build a module but want to build an application that also uses that module, let's go through the process of configuring an app.


II) xNFTs Module
----


### Messages

Messages are a great place to start when building a module because they define the actions that your application can make. Think of all the scenarios where a user would be able to update the state of the application in any way. These should be boiled down into basic interactions, similar to **CRUD** (Create, Read, Update, Delete).

Let's start with **Create**

**MsgXNFTTransfer**

Messages are types which live inside the `./xnfts/internal/types/` directory. There is already a `msg.go` file. We can use this `msg.go`.

```go=
package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
    "github.com/golang/protobuf/proto"
    
    "github.com/FreeFlixMedia/modules/nfts"
)

type NFTInput struct {
    PrimaryNFTID string `json:"primary_nft_id"`
    Recipient    string `json:"recipient"`
    AssetID      string `json:"asset_id"`
    
    LicensingFee  sdk.Coin `json:"licensing_fee"`
    RevenueShare  sdk.Dec  `json:"revenue_share"`
    TwitterHandle string   `json:"twitter_handle"`
}

type MsgXNFTTransfer struct {
    SourcePort    string         `json:"source_port"`
    SourceChannel string         `json:"source_channel"`
    DestHeight    uint64         `json:"dest_height"`
    Sender        sdk.AccAddress `json:"sender"`
    
    NFTInput
}

func NewMsgXNFTTransfer(sourcePort, sourceChannel string, height uint64, sender sdk.AccAddress,
    nftInput NFTInput) MsgXNFTTransfer {
    return MsgXNFTTransfer{
        SourcePort:    sourcePort,
        SourceChannel: sourceChannel,
        DestHeight:    height,
        Sender:        sender,
        NFTInput:      nftInput,
    }
}

var _ sdk.Msg = MsgXNFTTransfer{}

func (m *MsgXNFTTransfer) Reset() {
    *m = MsgXNFTTransfer{}
}

func (m *MsgXNFTTransfer) String() string {
    return proto.CompactTextString(m)
}

func (m MsgXNFTTransfer) ProtoMessage() {}

func (m MsgXNFTTransfer) Route() string {
    return RouterKey
}

func (m MsgXNFTTransfer) Type() string {
    return "msg_xnft_transfer"
}

func (m MsgXNFTTransfer) ValidateBasic() error {
    if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
        return sdkerrors.Wrap(err, "invalid source port ID")
    }
    if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
        return sdkerrors.Wrap(err, "invalid source channel ID")
    }
    
    if nfts.GetContextOfCurrentChain() == nfts.CoCoContext {
        if m.NFTInput.AssetID == "" {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
        } else if m.NFTInput.RevenueShare.IsZero() {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share is not allowed to be empty")
        } else if m.NFTInput.TwitterHandle == "" {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "handle name should not be empty")
        } else if !m.NFTInput.LicensingFee.IsValid() {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "licensing fee is invalid")
        }
    }
    
    if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext {
        if len(m.PrimaryNFTID) == 0 {
            return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nft id is empty")
        }
    }
    if m.Sender.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
    } else if m.NFTInput.Recipient == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Recipient should not be nil")
    }
    return nil
}

func (m MsgXNFTTransfer) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgXNFTTransfer) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{m.Sender}
}

```

Notice that all Messages in the app need to follow the `sdk.Msg` interface. The Message `struct` contains all the necessary information when creating a new xnft:

* `Sender`- Who initiates xnft transfer. This uses the `sdk.AccAddress` type which represents an account in the app.
* `SourcePort`- source chain port
* `SourceChannel`- source chain channel
* `DestHeight`- Destination chain height
* `NFTInput`:
    * `PrimaryNFTID`- Primary TweetNFT ID
    * `Recipient`- Primary nft receiver
    * `AssetID`- AssetID of the tweet nft.
    * `LicensingFee`- Amount to pay to get the license of an asset
    * `RevenueShare`- Amount of share we get when nft is used
    * `TwitterHandle`- Twitter handle of the user.


The `Msg` interface requires some other methods to be set, like validating the content of the `struct`, and confirming the msg was signed and submitted by the Creator.


### Codec

Once we have defined our messages, we need to describe to our encoder how they should be stored as bytes. To do this we edit the file located at `./xnfts/internal/types/codec.go`. By describing our types as follows they will work with our encoding library:

```go=
package types

import (
    "github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
    commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
)

var (
    amino     = codec.New()
    ModuleCdc = codec.NewHybridCodec(amino, cdctypes.NewInterfaceRegistry())
)

// RegisterCodec registers the IBC transfer types
func RegisterCodec(cdc *codec.Codec) {
    cdc.RegisterConcrete(MsgXNFTTransfer{}, "ibc/xnfts/MsgXNFTTransfer", nil)
    cdc.RegisterConcrete(BaseNFTPacket{}, "ibc/xnfts/BaseNFTPacket", nil)
    
    cdc.RegisterInterface((*XNFTs)(nil), nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
    registry.RegisterImplementations((*sdk.Msg)(nil), &MsgXNFTTransfer{})
}

func init() {
    RegisterCodec(amino)
    channel.RegisterCodec(amino)
    commitmenttypes.RegisterCodec(amino)
    amino.Seal()
}

```
### Alias

Now that we have these new message types, we'd like to make sure other parts of the module can access them. To do so we use the `./xnfts/alias.go` file. This imports the types from the nested `types` directory and makes them accessible at the module's top-level directory.

```go=
package xnfts

import (
    "github.com/FreeFlixMedia/modules/xnfts/internal/keeper"
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

type (
    Keeper = keeper.Keeper
    
    MsgXNFTTransfer = types.MsgXNFTTransfer
    BaseNFTPacket   = types.BaseNFTPacket
    
    PostCreationPacketAcknowledgement = types.PostCreationPacketAcknowledgement
    
    XNFTs = types.XNFTs
)

const (
    ModuleName   = types.ModuleName
    StoreKey     = types.StoreKey
    RouterKey    = types.RouterKey
    QuerierRoute = types.QuerierRoute
)

var (
    NewKeeper                  = keeper.NewKeeper
    RegisterCodec              = types.RegisterCodec
    RegisterInterfaces         = types.RegisterInterfaces
    AttributeValueCategory     = types.AttributeValueCategory
    AttributeKeyReceiver       = types.AttributeKeyReceiver
    EventTypeNFTPacketTransfer = types.EventTypeNFTPacketTransfer
)

```

It's great to have Messages, but we need somewhere to store the information they are sending. All persistent data related to this module should live in the module's `Keeper`.

Let's make a `Keeper` for our XNFTs Module next.


### Keeper

After using the scaffold command you should have a boilerplate Keeper at `./xnfts/internal/keeper/keeper.go`. It contains a basic keeper with references to basic functions like `Set`, `Get` and `Delete`.

Our keeper stores all our data for our module. Sometimes a module will import the keeper of another module. This will allow the state to be shared and modified across modules. Here we are dealing with the creation of nfts. Look at our completed Keeper and how to Set and Get are expanded:

```go=
package keeper

import (
    "fmt"
    
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/cosmos/cosmos-sdk/x/capability"
    channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
    channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
    ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
    "github.com/tendermint/tendermint/libs/log"
    
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

type Keeper struct {
    storeKey  sdk.StoreKey
    cdc       *codec.Codec
    nftKeeper types.NFTKeeper
    
    bankKeeper    types.BaseBankKeeper
    channelKeeper types.ChannelKeeper
    portKeeper    types.PortKeeper
    scopedKeeper  capability.ScopedKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, nftKeeper types.NFTKeeper, bankKeeper types.BaseBankKeeper, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
    scopedKeeper capability.ScopedKeeper) Keeper {
    
    return Keeper{
        
        storeKey:      key,
        cdc:           cdc,
        nftKeeper:     nftKeeper,
        bankKeeper:    bankKeeper,
        channelKeeper: channelKeeper,
        portKeeper:    portKeeper,
        scopedKeeper:  scopedKeeper,
    }
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
    return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.ModuleName))
}

func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
    chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
    if !ok {
        return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
    }
    
    return k.channelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
    chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(portID, channelID))
    if !ok {
        return sdkerrors.Wrapf(channel.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", ibctypes.ChannelCapabilityPath(portID, channelID))
    }
    
    return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

func (k Keeper) IsBounded(ctx sdk.Context, portID string) bool {
    _, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.PortPath(portID))
    return ok
}

func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
    store := ctx.KVStore(k.storeKey)
    
    store.Set([]byte(types.PortKey), []byte(portID))
    capa := k.portKeeper.BindPort(ctx, portID)
    return k.ClaimCapability(ctx, capa, ibctypes.PortPath(portID))
}

func (k Keeper) GetPort(ctx sdk.Context) string {
    store := ctx.KVStore(k.storeKey)
    return string(store.Get([]byte(types.PortKey)))
}

func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
    return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

```

### BaseNFTPacket

You may notice the reference to `types.BaseNFTPacket` throughout the `Keeper`. These is new struct defined in `./xnfts/types/packet.go` that contains all necessary information about different nft packets. You can create this file now and add the following:

```go=
package types

import (
    "encoding/json"
    "fmt"
    
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    
    "github.com/FreeFlixMedia/modules/nfts"
)

type XNFTs interface {
    GetBytes() []byte
}

type BaseNFTPacket struct {
    PrimaryNFTID    string `json:"primary_nftid"`
    PrimaryNFTOwner string `json:"primary_nft_owner"`
    
    SecondaryNFTID    string `json:"secondary_nftid"`
    SecondaryNFTOwner string `json:"secondary_nft_owner"`
    
    AssetID string `json:"asset_id"`
    
    License      bool     `json:"license"`
    LicensingFee sdk.Coin `json:"licensing_fee"`
    RevenueShare sdk.Dec  `json:"revenue_share"`
    
    TwitterHandle string `json:"twitter_handle"`
}

var _ XNFTs = BaseNFTPacket{}

func (nft *BaseNFTPacket) Reset() {
    *nft = BaseNFTPacket{}
}

func (nft BaseNFTPacket) ProtoMessage() {
}

func NewBaseNFTPacket(primaryNFTID, secondaryNFTID, primaryNFTOwner, secondaryNFTOwner string,
    assetID, twitterHandle string, license bool, fee sdk.Coin, share sdk.Dec) BaseNFTPacket {
    return BaseNFTPacket{
        PrimaryNFTID:      primaryNFTID,
        PrimaryNFTOwner:   primaryNFTOwner,
        SecondaryNFTID:    secondaryNFTID,
        SecondaryNFTOwner: secondaryNFTOwner,
        AssetID:           assetID,
        License:           license,
        LicensingFee:      fee,
        RevenueShare:      share,
        TwitterHandle:     twitterHandle,
    }
}

func (nft BaseNFTPacket) String() string {
    return fmt.Sprintf(`
PrimaryNFTID: %s
PrimaryNFTOwner: %s

SecondaryNFTID: %s
SecondaryNFTOwner: %s

AssetID: %s
License: %t

LicensingFee: %s
RevenueShare: %s
TwittterHandle: %s
`, nft.PrimaryNFTID, nft.PrimaryNFTOwner, nft.SecondaryNFTID, nft.SecondaryNFTOwner, nft.AssetID, nft.License,
        nft.LicensingFee, nft.RevenueShare, nft.TwitterHandle)
}

func (nft BaseNFTPacket) ValidateBasic() error {
    if nft.PrimaryNFTID == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nftd should be present")
    } else if nft.PrimaryNFTOwner == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nft owner is empty")
    } else if nft.SecondaryNFTID == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nft id is empty")
    } else if nft.SecondaryNFTOwner == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nft owner is empty")
    } else if nft.AssetID == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id is empty")
    } else if nft.RevenueShare.IsNil() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share is empty")
    } else if nft.TwitterHandle == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "handle name is empty")
    }
    return nil
}

func (nft BaseNFTPacket) GetBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(nft))
}

func (nft BaseNFTPacket) MarshalJSON() ([]byte, error) {
    type tmp BaseNFTPacket
    return json.Marshal(tmp(nft))
}

func (nft *BaseNFTPacket) UnmarshalJSON(bytes []byte) error {
    type tmp BaseNFTPacket
    var data tmp
    
    if err := json.Unmarshal(bytes, &data); err != nil {
        return err
    }
    
    *nft = BaseNFTPacket(data)
    return nil
}

func (nft BaseNFTPacket) ToBaseTweetNFT() *nfts.BaseTweetNFT {
    return &nfts.BaseTweetNFT{
        PrimaryNFTID:   nft.PrimaryNFTID,
        PrimaryOwner:   nft.PrimaryNFTOwner,
        SecondaryNFTID: nft.SecondaryNFTID,
        SecondaryOwner: nft.SecondaryNFTOwner,
        License:        nft.License,
        AssetID:        nft.AssetID,
        LicensingFee:   nft.LicensingFee,
        RevenueShare:   nft.RevenueShare,
        TwitterHandle:  nft.TwitterHandle,
    }
}

type PostCreationPacketAcknowledgement struct {
    Success bool   `json:"success" yaml:"success"`
    Error   string `json:"error" yaml:"error"`
}

func (ack PostCreationPacketAcknowledgement) GetBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}

```

You might also notice that each type has the `String` method. This allows us to render the struct as a string for rendering.

**ToBaseTweetNFT**

We use `ToBaseTweetNFT` to convert type from `BaseNFTPacket` to `BaseTweetNFT`

### Prefixes

You may notice the use of types.PortKey. These are defined in a file called `./xnfts/internal/types/key.go` and help us keep our `Keeper` organized. The `Keeper` is just a key-value store. That means that similar to an `Object` in javascript, all values are referenced under a key. To access a value, you need to know the key under which it is stored. This is a bit like a unique identifier (UID).

 You should add these to your `key.go` file so it looks as follows:

```go=
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
    ModuleName = "xnfts"
    Version    = "ics20-1"
    PortID     = ModuleName
    
    StoreKey     = ModuleName
    QuerierRoute = ModuleName
    RouterKey    = ModuleName
    PortKey      = "portID"
)

func GetHexAddressFromBech32String(addr string) sdk.AccAddress {
    addrs, _ := sdk.AccAddressFromBech32(addr)
    return addrs
}

```

### Handler

For a **Message** to reach a **Keeper**, it has to go through a **Handler**. This is where logic can be applied to either allow or deny a `Message` to succeed. It’s also where logic as to exactly how the state should change within the Keeper should take place. If you’re familiar with [Model View Controller](https://en.wikipedia.org/wiki/Model–view–controlle) (MVC) architecture, the `Keeper` is a bit like the **Model** and the `Handler` is a bit like the **Controller**. If you’re familiar with [React/Redux](https://en.wikipedia.org/wiki/React_(web_framework)) or [Vue/Vuex](https://en.wikipedia.org/wiki/Vue.js) architecture, the `Keeper` is a bit like the **Reducer/Store** and the `Handler` is a bit like **Actions**.


Our Handler will go in `./xnfts/handler.go` and will follow the suggestions outlined in the boilerplate. We will create handler functions for `Message` type, `MsgXNFTTransfer` until the file looks as follows:

```go=
package xnfts

import (
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
    
    "github.com/FreeFlixMedia/modules/nfts"
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        ctx = ctx.WithEventManager(sdk.NewEventManager())
        
        switch msg := msg.(type) {
        case MsgXNFTTransfer:
            return handleMsgXNFTTransfer(ctx, k, msg)
        
        default:
            return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 xnft message type: %T", msg)
            
        }
    }
}

func handleMsgXNFTTransfer(ctx sdk.Context, k Keeper, msg MsgXNFTTransfer) (*sdk.Result, error) {
    var packet BaseNFTPacket
    
    if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext {
        nft, found := k.GetTweetNFTByID(ctx, msg.PrimaryNFTID)
        if !found {
            return nil, sdkerrors.Wrap(nfts.ErrNFTNotFound, "")
        }
        
        if !nft.License {
            return nil, sdkerrors.Wrap(nfts.ErrInvalidLicense, fmt.Sprintf("unable to transfer %s", nft.PrimaryNFTID))
        }
        
        if !msg.Sender.Equals(types.GetHexAddressFromBech32String(nft.PrimaryOwner)) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "")
        }
        
        packet.PrimaryNFTID = nft.PrimaryNFTID
        packet.PrimaryNFTOwner = nft.PrimaryOwner
        packet.License = nft.License
        packet.AssetID = nft.AssetID
        packet.RevenueShare = nft.RevenueShare
        packet.LicensingFee = nft.LicensingFee
        packet.SecondaryNFTOwner = msg.Recipient
        packet.TwitterHandle = nft.TwitterHandle
        
    } else if nfts.GetContextOfCurrentChain() == nfts.CoCoContext {
        
        count := k.GetGlobalTweetCount(ctx)
        sNFTID := nfts.GetSecondaryNFTID(count)
        
        packet.PrimaryNFTOwner = msg.Recipient
        packet.License = true
        packet.AssetID = msg.AssetID
        packet.RevenueShare = msg.RevenueShare
        packet.LicensingFee = msg.LicensingFee
        packet.SecondaryNFTID = sNFTID
        packet.SecondaryNFTOwner = msg.Sender.String()
        packet.TwitterHandle = msg.TwitterHandle
        
        k.MintTweetNFT(ctx, *packet.ToBaseTweetNFT())
        k.SetTweetIDToAccount(ctx, msg.Sender, sNFTID)
        k.SetGlobalTweetCount(ctx, count+1)
    }
    
    if err := k.XTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, packet.GetBytes()); err != nil {
        return nil, err
    }
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            sdk.EventTypeMessage,
            sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
            sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
            sdk.NewAttribute(AttributeKeyReceiver, msg.Recipient),
        ),
    )
    return &sdk.Result{
        Events: ctx.EventManager().Events().ToABCIEvents(),
    }, nil
}

func handleXNFTRecvPacket(ctx sdk.Context, k Keeper, packet channeltypes.Packet) (*sdk.Result, error) {
    
    var nftData BaseNFTPacket
    if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &nftData); err != nil {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
    }
    
    acknowledgement := PostCreationPacketAcknowledgement{
        Success: true,
        Error:   "",
    }
    
    if err := k.OnRecvNFTPacket(ctx, nftData, packet); err != nil {
        acknowledgement = PostCreationPacketAcknowledgement{
            Success: false,
            Error:   err.Error(),
        }
    }
    
    if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
        return nil, err
    }
    
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            EventTypeNFTPacketTransfer,
            sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
        ),
    )
    
    return &sdk.Result{
        Events: ctx.EventManager().Events().ToABCIEvents(),
    }, nil
}

```

### Events

At the end of each handler is an EventManager which will create logs within the transaction that reveals information about what occurred during the handling of this message. This is useful for client-side software that wants to know exactly what happened as a result of this state transition. These Events use a series of pre-defined types that can be found in `./xnfts/internal/types/events.go` and look as follows:

```go=
package types

import (
    "fmt"
    
    ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

var (
    EventTypeNFTPacketTransfer = "nft_packet_transfer"
    AttributeKeyReceiver       = "receiver"
    AttributeValueCategory     = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)

```

Now that we have all the necessary pieces for updating state (`Message`, `Handler`, `Keeper`) we might want to consider ways in which we can query state. This is typically done via a REST endpoint and/or a CLI. Both of those clients interact with a part of the app which queries state, called the `Querier`.

### relay.go

To relay the `NFTPacket`, we will access the channel created between two chains and the portId of the destination chain. Then we will create an outgoing packet that is `createOutgoingPacket`. After receiving the packet in destination chain we will perform the logic that defined in the function `handleXNFTRecvPacket` in `handler.go` file

```go=
package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
    ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
    
    "github.com/FreeFlixMedia/modules/nfts"
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

const (
    DefaultPacketTimeoutHeight = 1000
    
    DefaultPacketTimeoutTimestamp = 0
)

func (k Keeper) XTransfer(
    ctx sdk.Context,
    sourcePort, sourceChannel string,
    destHeight uint64,
    packetData []byte,
) error {
    
    sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
    if !found {
        return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
    }
    
    destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
    destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()
    
    sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
    if !found {
        return channeltypes.ErrSequenceSendNotFound
    }
    
    return k.createOutgoingPacket(ctx, sequence, sourcePort, sourceChannel, destinationPort, destinationChannel, destHeight, packetData)
}

func (k Keeper) createOutgoingPacket(
    ctx sdk.Context,
    seq uint64,
    sourcePort, sourceChannel string,
    destinationPort, destinationChannel string,
    destHeight uint64,
    data []byte,
) error {
    
    channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
    if !ok {
        return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
    }
    
    packet := channeltypes.NewPacket(
        data,
        seq,
        sourcePort,
        sourceChannel,
        destinationPort,
        destinationChannel,
        destHeight+DefaultPacketTimeoutHeight, // TODO : DestHeight need to be updated with src header.height
        DefaultPacketTimeoutTimestamp,
    )
    
    return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) OnRecvNFTPacket(ctx sdk.Context, data types.BaseNFTPacket, packet channeltypes.Packet) error {
    
    if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext && len(data.PrimaryNFTID) == 0 {
        addr, err := sdk.AccAddressFromBech32(data.PrimaryNFTOwner)
        if err != nil {
            return err
        }
        
        _, err = k.bankKeeper.AddCoins(ctx, addr, sdk.Coins{data.LicensingFee})
        if err != nil {
            return err
        }
        
        count := k.nftKeeper.GetGlobalTweetCount(ctx)
        primaryNFTID := nfts.GetPrimaryNFTID(count)
        data.PrimaryNFTID = primaryNFTID
        
        k.nftKeeper.MintTweetNFT(ctx, *data.ToBaseTweetNFT())
        k.SetTweetIDToAccount(ctx, addr, primaryNFTID)
        k.SetGlobalTweetCount(ctx, count+1)
        
        if err := k.XTransfer(ctx, packet.DestinationPort, packet.DestinationChannel, packet.TimeoutHeight, data.GetBytes()); err != nil {
            return err
        }
        
    }
    if nfts.GetContextOfCurrentChain() == nfts.CoCoContext && len(data.SecondaryNFTID) == 0 {
        addr, err := sdk.AccAddressFromBech32(data.SecondaryNFTOwner)
        if err != nil {
            return err
        }
        
        count := k.nftKeeper.GetGlobalTweetCount(ctx)
        secondaryNFTID := nfts.GetSecondaryNFTID(count)
        data.SecondaryNFTID = secondaryNFTID
        
        k.nftKeeper.MintTweetNFT(ctx, *data.ToBaseTweetNFT())
        k.nftKeeper.SetTweetIDToAccount(ctx, addr, secondaryNFTID)
        k.nftKeeper.SetGlobalTweetCount(ctx, count+1)
        
    }
    
    ctx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            nfts.EventTypeMsgMintTweetNFT,
            sdk.NewAttribute(sdk.AttributeKeySender, data.PrimaryNFTOwner),
            sdk.NewAttribute(types.AttributeKeyReceiver, data.SecondaryNFTOwner),
            sdk.NewAttribute(nfts.AttributePrimaryNFTID, data.PrimaryNFTID),
            sdk.NewAttribute(nfts.AttributeSecondaryNFTID, data.SecondaryNFTID),
        ),
    })
    return nil
}

```
we have all of the basic actions of our module created, we want to make them accessible. We can do this with a CLI client and a REST client. For this tutorial we will just be creating a CLI client. If you are interested in what goes into making a REST client.

Let’s take a look at what goes into making a CLI.

### CLI

A Command Line Interface (CLI) will help us interact with our app once it is running on a machine somewhere. Each Module has its own namespace within the CLI that gives it the ability to create and sign Messages destined to be handled by that module. It also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI will let you do things like generate keys for a new account or check the status of an interaction you already had with the application.


The CLI for our module is broken into two files called tx.go and query.go which are located in `./xnfts/client/cli/.` One file is for making transactions that contain messages which will ultimately update our state. The other is for making queries which will give us the ability to read information from our state. Both files utilize the [Cobra](https://github.com/spf13/cobra) library.

**Transactions**

The `tx.go` file contains `GetTxCmd` which is a standard method within the Cosmos SDK. It is referenced later in the `module.go` file which describes exactly which attributes a module has. This makes it easier to incorporate different modules for different reasons at the level of the actual application. After all, we are focusing on a module at this point, but later we will create an application that utilizes this module as well as other modules that are already available within the Cosmos SDK.

Inside `GetTxCmd` we create a new module-specific command and call is `xnfts`. Within this command we add a sub-command for each Message type we’ve defined:
* `GetXNFTTxCmd`

Each function takes parameters from the **Cobra** CLI tool to create a new msg, sign it and submit it to the application to be processed. These functions should go into the `tx.go` file and look as follows:

```go=
package cli

import (
    "bufio"
    "strconv"
    
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

// GetTxCmd returns the transaction commands for IBC fungible token transfer
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
    ics20XNFTTransferTxCmd := &cobra.Command{
        Use:   "xnfts",
        Short: "IBC nft  transfer transaction subcommands",
    }
    
    ics20XNFTTransferTxCmd.AddCommand(flags.PostCommands(
        GetXNFTTxCmd(cdc),
    )...)
    
    return ics20XNFTTransferTxCmd
}

// GetXNFTTxCmd returns the command to create a NewMsgTransfer transaction
func GetXNFTTxCmd(cdc *codec.Codec) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "nft-transfer [src-port] [src-channel] [dest-height] [recipient]",
        Short: "Transfer non fungible token through IBC",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
            cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
            
            sender := cliCtx.GetFromAddress()
            srcPort := args[0]
            srcChannel := args[1]
            destHeight, err := strconv.Atoi(args[2])
            if err != nil {
                return err
            }
            
            var fee sdk.Coin
            var share sdk.Dec
            var assteID, handle string
            var data types.NFTInput
            
            licenceFee := viper.GetString(FlagLicensingFee)
            if licenceFee != "" {
                fee, err = sdk.ParseCoin(licenceFee)
                if err != nil {
                    return err
                }
            }
            
            shareStr := viper.GetString(FlagRevenueShare)
            if shareStr != "" {
                share, err = sdk.NewDecFromStr(shareStr)
                if err != nil {
                    return err
                }
                
            }
            
            assteID = viper.GetString(FlagAssetID)
            handle = viper.GetString(FlagTwitterHandle)
            
            data = types.NFTInput{
                PrimaryNFTID:  viper.GetString(FlagPrimaryNFTID),
                Recipient:     args[3],
                AssetID:       assteID,
                LicensingFee:  fee,
                RevenueShare:  share,
                TwitterHandle: handle,
            }
            
            msg := types.NewMsgXNFTTransfer(srcPort, srcChannel, uint64(destHeight), sender, data)
            if err := msg.ValidateBasic(); err != nil {
                return err
            }
            
            return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
    
    cmd.Flags().String(FlagPrimaryNFTID, "", "Primary NFT id")
    cmd.Flags().String(FlagRevenueShare, "0", "Revenue share")
    cmd.Flags().String(FlagLicensingFee, "0coco", "Licenese fee")
    cmd.Flags().String(FlagAssetID, "", "AssetID")
    cmd.Flags().String(FlagTwitterHandle, "", "Twitter Handle")
    return cmd
}

```

### Module


Our `scaffold` tool has done most of the work for us in generating our `module.go` file inside `./xnfts/.` One way that our module is different than the simplest form of a module, is that it uses it’s own `Keeper`. The only real changes needed are under the `AppModule` and `NewAppModule`. The file should look as follows afterward:

```go=
package xnfts

import (
    "encoding/json"
    
    "github.com/gorilla/mux"
    "github.com/spf13/cobra"
    
    "github.com/FreeFlixMedia/modules/xnfts/client/cli"
    
    abci "github.com/tendermint/tendermint/abci/types"
    
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/cosmos/cosmos-sdk/types/module"
    "github.com/cosmos/cosmos-sdk/x/capability"
    channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
    channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
    port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
    porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
    ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
    
    "github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

var (
    _ module.AppModule      = AppModule{}
    _ port.IBCModule        = AppModule{}
    _ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
    RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
    return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ json.RawMessage) error {
    return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {

}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
    return nil
}

func (AppModuleBasic) RegisterInterfaceTypes(registry cdctypes.InterfaceRegistry) {
    RegisterInterfaces(registry)
}

type AppModule struct {
    AppModuleBasic
    keeper Keeper
}

func NewAppModule(k Keeper) AppModule {
    return AppModule{
        keeper: k,
    }
}

func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
    // TODO
}

func (AppModule) Route() string {
    return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
    return NewHandler(am.keeper)
}

func (AppModule) QuerierRoute() string {
    return QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
    return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
    var genesisState types.GenesisState
    cdc.MustUnmarshalJSON(data, &genesisState)
    
    InitGenesis(ctx, am.keeper, genesisState)
    return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
    gs := ExportGenesis(ctx, am.keeper)
    return cdc.MustMarshalJSON(gs)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {

}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
    return []abci.ValidatorUpdate{}
}

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
    ctx sdk.Context,
    order ibctypes.Order,
    connectionHops []string,
    portID string,
    channelID string,
    chanCap *capability.Capability,
    counterparty channeltypes.Counterparty,
    version string,
) error {
    // TODO: Enforce ordering, currently relayers use ORDERED channels
    
    // Require portID is the portID transfer module is bound to
    boundPort := am.keeper.GetPort(ctx)
    if boundPort != portID {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
    }
    
    if version != types.Version {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
    }
    
    // Claim channel capability passed back by IBC module
    if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
        return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
    }
    
    // TODO: escrow
    return nil
}

func (am AppModule) OnChanOpenTry(
    ctx sdk.Context,
    order ibctypes.Order,
    connectionHops []string,
    portID,
    channelID string,
    chanCap *capability.Capability,
    counterparty channeltypes.Counterparty,
    version,
    counterpartyVersion string,
) error {
    // TODO: Enforce ordering, currently relayers use ORDERED channels
    
    // Require portID is the portID transfer module is bound to
    boundPort := am.keeper.GetPort(ctx)
    if boundPort != portID {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
    }
    
    if version != types.Version {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
    }
    
    if counterpartyVersion != types.Version {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
    }
    
    // Claim channel capability passed back by IBC module
    if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
        return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
    }
    
    // TODO: escrow
    return nil
}

func (am AppModule) OnChanOpenAck(
    ctx sdk.Context,
    portID,
    channelID string,
    counterpartyVersion string,
) error {
    if counterpartyVersion != types.Version {
        return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
    }
    return nil
}

func (am AppModule) OnChanOpenConfirm(
    ctx sdk.Context,
    portID,
    channelID string,
) error {
    return nil
}

func (am AppModule) OnChanCloseInit(
    ctx sdk.Context,
    portID,
    channelID string,
) error {
    // Disallow user-initiated channel closing for transfer channels
    return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

func (am AppModule) OnChanCloseConfirm(
    ctx sdk.Context,
    portID,
    channelID string,
) error {
    return nil
}

func (am AppModule) OnRecvPacket(
    ctx sdk.Context,
    packet channeltypes.Packet,
) (*sdk.Result, error) {
    var data XNFTs
    
    if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 xnft packet data: %s", err.Error())
    }
    
    switch data := data.(type) {
    case BaseNFTPacket:
        return handleXNFTRecvPacket(ctx, am.keeper, packet)
    default:
        return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 transfer message type: %T", data)
        
    }
    
}

func (am AppModule) OnAcknowledgementPacket(
    ctx sdk.Context,
    packet channeltypes.Packet,
    acknowledgement []byte,
) (*sdk.Result, error) {
    
    return &sdk.Result{
        Events: ctx.EventManager().Events().ToABCIEvents(),
    }, nil
}

func (am AppModule) OnTimeoutPacket(
    ctx sdk.Context,
    packet channeltypes.Packet,
) (*sdk.Result, error) {
    
    return &sdk.Result{
        Events: ctx.EventManager().Events().ToABCIEvents(),
    }, nil
}

```

Congratulations you have completed the `xnfts` module!

This module is now able to be incorporated into any Cosmos SDK application.

Since we don’t want to just build a module but want to build an application that also uses that module, let’s go through the process of configuring an app.


