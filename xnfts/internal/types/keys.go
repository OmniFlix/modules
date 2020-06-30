package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "xnfts"
	Version    = "ics20-1"
	PortID     = "xnfts"
	
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	PortKey      = "portID"
	
	QueryParams       = "params"
	TransferPortID    = "transfer"
	DefaultParamspace = ModuleName
)

var (
	LastVisitedTime = []byte{0x01}
	ParamKey        = []byte{0x02}
)

func GetParamKey() []byte {
	return ParamKey
}

func GetLastVisitedKey() []byte {
	return LastVisitedTime
}

func GetHexAddressFromBech32String(addr string) sdk.AccAddress {
	addrs, _ := sdk.AccAddressFromBech32(addr)
	return addrs
}
