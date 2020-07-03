
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