package xnfts

import (
	"github.com/FreeFlixMedia/modules/xnfts/internal/keeper"
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

const (
	ModuleName                             = types.ModuleName
	RouterKey                              = types.RouterKey
	QuerierRoute                           = types.QuerierRoute
	Version                                = types.Version
	StoreKey                               = types.StoreKey
	XNFTPortID                             = types.PortID
	AttributeKeyReceiver                   = types.AttributeKeyReceiver
	EventTypeNFTPacketTransfer             = types.EventTypeNFTPacketTransfer
	EventTypeTokenDistribution             = types.EventTypeTokenDistribution
	EventTypePayLicensingFeeAndNFTTransfer = types.EventTypePayLicensingFeeAndNFTTransfer
	AttributeKeyAckSuccess                 = types.AttributeKeyAckSuccess
	AttributeKeyAckError                   = types.AttributeKeyAckError
	AttributeKeyNFTChannel                 = types.AttributeKeyNFTChannel
	AttributeKeyDestHeight                 = types.AttributeKeyDestHeight
	EventTypeSetParams                     = types.EventTypeSetParams
)

type (
	Keeper                              = keeper.Keeper
	GenesisState                        = types.GenesisState
	XNFTs                               = types.XNFTs
	Params                              = types.Params
	NFTInput                            = types.NFTInput
	MsgXNFTTransfer                     = types.MsgXNFTTransfer
	MsgSlotBooking                      = types.MsgSlotBooking
	MsgSetParams                        = types.MsgSetParams
	MsgPayLicensingFee                  = types.MsgPayLicensingFee
	MsgDistributeFunds                  = types.MsgDistributeFunds
	ChannelKeeper                       = types.ChannelKeeper
	PostCreationPacketAcknowledgement   = types.PostCreationPacketAcknowledgement
	PacketSlotBooking                   = types.PacketSlotBooking
	BaseNFTPacket                       = types.BaseNFTPacket
	PacketTokenDistribution             = types.PacketTokenDistribution
	PacketPayLicensingFeeAndNFTTransfer = types.PacketPayLicensingFeeAndNFTTransfer
)

var (
	NewKeeper                     = keeper.NewKeeper
	NewParams                     = types.NewParams
	NewQuerier                    = keeper.NewQuerier
	DefaultParamspace             = types.DefaultParamspace
	ModuleCdc                     = types.ModuleCdc
	RegisterCodec                 = types.RegisterCodec
	DefaultGenesis                = types.DefaultGenesis
	RegisterInterfaces            = types.RegisterInterfaces
	GetHexAddressFromBech32String = types.GetHexAddressFromBech32String
	NewMsgXNFTTransfer            = types.NewMsgXNFTTransfer
	NewMsgSetParams               = types.NewMsgSetParams
	NewMsgSlotBooking             = types.NewMsgSlotBooking
	NewMsgPayLicensingFee         = types.NewMsgPayLicensingFee
	NewMsgDistributeFunds         = types.NewMsgDistributeFunds
	ErrNFTNotFound                = types.ErrNFTNotFound
	AttributeValueCategory        = types.AttributeValueCategory
)
