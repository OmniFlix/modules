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
	
	FreeFlixNFTPrefix  = "ffmttweetnft"
	CoCoNFTPrefix      = "cocotweetnft"
	LiveStreamIDPrefix = "cocolivestream"
	FreeFlixAdNFT      = "ffmtadnft"
	
	FreeFlixContext = "freeflix"
	CoCoContext     = "coco"
)

var (
	GlobalTweetCountPrefix   = []byte{0x01}
	GlobalAdsCountPrefix     = []byte{0x02}
	TweetNFTPrefix           = []byte{0x03}
	AdNFTPrefix              = []byte{0x04}
	TweetAccountPrefix       = []byte{0x05}
	AdsAccountPrefix         = []byte{0x06}
	TwitterHandlePrefix      = []byte{0x07}
	DNFTPrefix               = []byte{0x08}
	LiveStreamPrefix         = []byte{0x09}
	GlobalLiveStreamCountKey = []byte{0x11}
	AclKey                   = []byte{0x12}
	AuthorisedHandlersKey    = []byte{0x13}
)

func GetGlobalTweetCountKey() []byte {
	return GlobalTweetCountPrefix
}

func GetGlobalAdsCountKey() []byte {
	return GlobalAdsCountPrefix
}

func GetTweetNFTKey(id []byte) []byte {
	return append(TweetNFTPrefix, id...)
}

func GetAdsNFTKey(id []byte) []byte {
	return append(AdNFTPrefix, id...)
}

func GetTweetsCountOfAddressKey(addr []byte) []byte {
	return append(TweetAccountPrefix, addr...)
}

func GetAdsCountOfAddressKey(addr []byte) []byte {
	return append(AdsAccountPrefix, addr...)
}

func GetPrimaryNFTID(count uint64) string {
	return FreeFlixNFTPrefix + strconv.Itoa(int(count))
}

func GetSecondaryNFTID(count uint64) string {
	return CoCoNFTPrefix + strconv.Itoa(int(count))
}

func GetLiveStreamID(count uint64) string {
	return LiveStreamIDPrefix + strconv.Itoa(int(count))
}

func GetAdNFTID(count uint64) string {
	return FreeFlixAdNFT + strconv.Itoa(int(count))
}

func GetContextOfCurrentChain() string {
	config := sdk.GetConfig()
	return config.GetBech32AccountAddrPrefix()
}

func GetDnftStoreKey(slotTime, programTime []byte) []byte {
	return append(GetDNFTSlotTimeKey(slotTime), programTime...)
}

func GetDNFTSlotTimeKey(slotTime []byte) []byte {
	return append(DNFTPrefix, slotTime...)
}

func GetLiveStreamKey(id []byte) []byte {
	return append(LiveStreamPrefix, id...)
}

func GetGlobalLiveStreamCountKey() []byte {
	return GlobalLiveStreamCountKey
}

func GetDNFTTimeKey(timeStamp []byte) []byte {
	return append(DNFTPrefix, timeStamp...)
}

func GetHandlerKey(handle []byte) []byte {
	return append(TwitterHandlePrefix, handle...)
}

func GetAclKey() []byte {
	return AclKey
}

func GetAuthorisedHandlersKey() []byte {
	return AuthorisedHandlersKey
}
