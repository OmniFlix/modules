
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
