package nfts

import (
	"github.com/FreeFlixMedia/modules/nfts/internal/keeper"
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

const (
	CoCoContext     = types.CoCoContext
	FreeFlixContext = types.FreeFlixContext
	
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	
	MsgMintTweetNFT = types.MsgMintTweetNFT
	BaseTweetNFT    = types.BaseTweetNFT
)

var (
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier
	GetContextOfCurrentChain = types.GetContextOfCurrentChain
	GetPrimaryNFTID          = types.GetPrimaryNFTID
	GetSecondaryNFTID        = types.GetSecondaryNFTID
	
	EventTypeMsgMintTweetNFT = types.EventTypeMsgMintTweetNFT
	
	AttributePrimaryNFTID   = types.AttributePrimaryNFTID
	AttributeSecondaryNFTID = types.AttributeSecondaryNFTID
	AttributeAssetID        = types.AttributeAssetID
	AttributeTwitterHandle  = types.AttributeTwitterHandle
	
	ErrAssetIDAlreadyExist = types.ErrAssetIDAlreadyExist
	ErrInvalidLicense      = types.ErrInvalidLicense
	ErrParamsNotFound      = types.ErrParamsNotFound
	ErrNFTNotFound         = types.ErrNFTNotFound
)
