
### Alias

Now that we have these new message types, we'd like to make sure other parts of the module can access them. To do so we use the `./xnfts/alias.go` file. This imports the types from the nested `types` directory and makes them accessible at the module's top-level directory.

```go=
package xnfts

import (
	"github.com/FreeFlixMedia/modules/xnfts/internal/keeper"
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

type (
	Keeper                              = keeper.Keeper
	BaseNFTPacket                       = types.BaseNFTPacket
	NFTInput                            = types.NFTInput
	XNFTs                               = types.XNFTs
	MsgXNFTTransfer                     = types.MsgXNFTTransfer
	MsgPayLicensingFee                  = types.MsgPayLicensingFee
	PostCreationPacketAcknowledgement   = types.PostCreationPacketAcknowledgement
	PacketPayLicensingFeeAndNFTTransfer = types.PacketPayLicensingFeeAndNFTTransfer
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper                              = keeper.NewKeeper
	RegisterCodec                          = types.RegisterCodec
	RegisterInterfaces                     = types.RegisterInterfaces
	NewMsgXNFTTransfer                     = types.NewMsgXNFTTransfer
	GetHexAddressFromBech32String          = types.GetHexAddressFromBech32String
	AttributeValueCategory                 = types.AttributeValueCategory
	AttributeKeyReceiver                   = types.AttributeKeyReceiver
	EventTypeNFTPacketTransfer             = types.EventTypeNFTPacketTransfer
	EventTypePayLicensingFeeAndNFTTransfer = types.EventTypePayLicensingFeeAndNFTTransfer
)

```

It's great to have Messages, but we need somewhere to store the information they are sending. All persistent data related to this module should live in the module's `Keeper`.

Let's make a `Keeper` for our XNFTs Module next.