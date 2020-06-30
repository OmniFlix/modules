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
