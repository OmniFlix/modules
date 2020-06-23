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
