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
