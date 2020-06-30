package nfts

import (
	"github.com/FreeFlixMedia/modules/nfts/internal/keeper"
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

type (
	BaseTweetNFT             = types.BaseTweetNFT
	BaseDNFT                 = types.BaseDNFT
	BaseLiveStream           = types.BaseLiveStream
	BaseAdNFT                = types.BaseAdNFT
	TwitterAccountInfo       = types.TwitterAccountInfo
	AllowedHandles           = types.AllowedHandles
	Keeper                   = keeper.Keeper
	GenesisState             = types.GenesisState
	MsgCreateInitialTweetNFT = types.MsgCreateInitialTweetNFT
	MsgMintTweetNFT          = types.MsgMintTweetNFT
	MsgLiveStream            = types.MsgLiveStream
	MsgUpdateLiveStream      = types.MsgUpdateLiveStream
	MsgBookSlot              = types.MsgBookSlot
	MsgCreateAdNFT           = types.MsgCreateAdNFT
	MsgClaimTwitterAccount   = types.MsgClaimTwitterAccount
	MsgUpdateHandlersInfo    = types.MsgUpdateHandlersInfo
	MsgUpdateAccessList      = types.MsgUpdateAccessList
)

const (
	ModuleName      = types.ModuleName
	QuerierRoute    = types.QuerierRoute
	StoreKey        = types.StoreKey
	RouteKey        = types.RouterKey
	FreeFlixContext = types.FreeFlixContext
	CoCoContext     = types.CoCoContext
	TypeLocal       = types.TypeLocal
	TypeIBC         = types.TypeIBC
	StatusActive    = types.StatusActive
	StatusInActive  = types.StatusInActive
)

var (
	NewKeeper                               = keeper.NewKeeper
	NewQuerier                              = keeper.NewQuerier
	GetPrimaryNFTID                         = types.GetPrimaryNFTID
	GetSecondaryNFTID                       = types.GetSecondaryNFTID
	GetContextOfCurrentChain                = types.GetContextOfCurrentChain
	GetLiveStreamID                         = types.GetLiveStreamID
	GetAdNFTID                              = types.GetAdNFTID
	GetAdsCountOfAddressKey                 = types.GetAdsCountOfAddressKey
	GenerateDNFTID                          = types.GenerateDNFTID
	GetTimeSlotFromDNFTID                   = types.GetTimeSlotFromDNFTID
	NewMsgCreateInitialTweetNFT             = types.NewMsgCreateInitialTweetNFT
	NewMsgMintNFT                           = types.NewMsgMintNFT
	NewMsgLiveStream                        = types.NewMsgLiveStream
	NewUpdateLiveStream                     = types.NewUpdateLiveStream
	NewMsgBookSlot                          = types.NewMsgBookSlot
	NewMsgCreateAdNFT                       = types.NewMsgCreateAdNFT
	NewMsgClaimTwitterAccount               = types.NewMsgClaimTwitterAccount
	NewMsgUpdateHandlerInfo                 = types.NewMsgUpdateHandlerInfo
	ErrNFTNotFound                          = types.ErrNFTNotFound
	ErrAssetIDAlreadyExist                  = types.ErrAssetIDAlreadyExist
	GetBech32StringOfSlotAndProgammeDetails = types.GetBech32StringOfSlotAndProgrammeDetails
	EventTypeDistribution                   = types.EventTypeDistribution
	EventTypeMsgBookSlot                    = types.EventTypeMsgBookSlot
	AttributeDNFTID                         = types.AttributeDNFTID
	AttributeLiveStreamID                   = types.AttributeLiveStreamID
)
