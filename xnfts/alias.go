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
