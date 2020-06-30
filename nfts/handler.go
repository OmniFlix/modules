package nfts

import (
	"fmt"
	"reflect"
	"strings"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		
		switch msg := msg.(type) {
		case MsgCreateInitialTweetNFT:
			return handleMsgCreateInitialTweetNFT(ctx, keeper, msg)
		case MsgMintTweetNFT:
			return handleMsgMintTweetNFT(ctx, keeper, msg)
		case MsgBookSlot:
			return handleMsgBookSlot(ctx, keeper, msg)
		case MsgLiveStream:
			return handleMsgLiveStreamCreation(ctx, keeper, msg)
		case MsgUpdateLiveStream:
			return handleMsgUpdateLiveStream(ctx, keeper, msg)
		case MsgCreateAdNFT:
			return handleMsgCreateAdNFT(ctx, keeper, msg)
		case MsgClaimTwitterAccount:
			return handleMsgClaimTwitterAccount(ctx, keeper, msg)
		case MsgUpdateAccessList:
			return handleMsgUpdateAccessList(ctx, keeper, msg)
		case MsgUpdateHandlersInfo:
			return handleMsgUpdateHandlers(ctx, keeper, msg)
		
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized NFT message type: %T", msg)
		}
	}
}

func handleMsgCreateInitialTweetNFT(ctx sdk.Context, keeper Keeper, msg MsgCreateInitialTweetNFT) (*sdk.Result, error) {
	count := keeper.GetGlobalTweetCount(ctx)
	
	_, found := keeper.GetTwitterHandleInfo(ctx, msg.TwitterHandle)
	if found {
		return nil, sdkerrors.Wrap(types.ErrAddressAlreadyExist, "handle already registered")
	}
	
	authorisedHandlers := keeper.GetAuthorisedHandlerInfo(ctx)
	_, found = authorisedHandlers.Find(msg.TwitterHandle)
	if !found {
		authorisedHandlers.Handles = append(authorisedHandlers.Handles, msg.TwitterHandle)
	}
	
	id := GetPrimaryNFTID(count)
	tweetNFT := BaseTweetNFT{
		PrimaryNFTID:   id,
		PrimaryOwner:   msg.Sender.String(),
		SecondaryNFTID: "",
		SecondaryOwner: "",
		License:        false,
		AssetID:        msg.AssetID,
		LicensingFee:   sdk.Coin{},
		RevenueShare:   sdk.Dec{},
		TwitterHandle:  msg.TwitterHandle,
	}
	
	keeper.MintTweetNFT(ctx, tweetNFT)
	keeper.SetTweetIDToAccount(ctx, msg.Sender, tweetNFT.PrimaryNFTID)
	keeper.SetGlobalTweetCount(ctx, count+1)
	
	info := TwitterAccountInfo{
		Owner:       msg.Sender,
		Handle:      msg.TwitterHandle,
		ClaimStatus: true,
	}
	keeper.SetTwitterHandlerInfo(ctx, info)
	keeper.UpdaterHandlerInfo(ctx, authorisedHandlers)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMsgCreateInitialNFT,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributePrimaryNFTID, tweetNFT.PrimaryNFTID),
			sdk.NewAttribute(types.AttributeAssetID, tweetNFT.AssetID),
			sdk.NewAttribute(types.AttributeTwitterHandle, tweetNFT.TwitterHandle),
		),
	)
	
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	
}

func handleMsgMintTweetNFT(ctx sdk.Context, keeper Keeper, msg MsgMintTweetNFT) (*sdk.Result, error) {
	
	_, found := keeper.GetTwitterHandleInfo(ctx, msg.TwitterHandle)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "twitter account is not initialized")
	}
	
	nfts := keeper.GetTweetsOfAccount(ctx, msg.Sender)
	
	for _, nft := range nfts {
		if strings.EqualFold(nft.AssetID, msg.AssetID) {
			return nil, sdkerrors.Wrap(types.ErrAssetIDAlreadyExist, "")
		}
	}
	
	count := keeper.GetGlobalTweetCount(ctx)
	id := GetPrimaryNFTID(count)
	tweetNFT := BaseTweetNFT{
		PrimaryNFTID:   id,
		PrimaryOwner:   msg.Sender.String(),
		SecondaryNFTID: "",
		SecondaryOwner: "",
		License:        msg.License,
		AssetID:        msg.AssetID,
		LicensingFee:   msg.LicensingFee,
		RevenueShare:   msg.RevenueShare,
		TwitterHandle:  msg.TwitterHandle,
	}
	
	keeper.MintTweetNFT(ctx, tweetNFT)
	keeper.SetTweetIDToAccount(ctx, msg.Sender, tweetNFT.PrimaryNFTID)
	keeper.SetGlobalTweetCount(ctx, count+1)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMsgMintTweetNFT,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributePrimaryNFTID, tweetNFT.PrimaryNFTID),
			sdk.NewAttribute(types.AttributeAssetID, tweetNFT.AssetID),
			sdk.NewAttribute(types.AttributeTwitterHandle, tweetNFT.TwitterHandle),
		),
	)
	
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	
}

// BookSlot

func handleMsgBookSlot(ctx sdk.Context, keeper Keeper, msg MsgBookSlot) (*sdk.Result, error) {
	
	sNFT, found := keeper.GetTweetNFTByID(ctx, msg.SecondaryNFTID)
	if !found {
		return nil, sdkerrors.Wrap(ErrNFTNotFound, "secondary nfts not found")
	}
	
	addr, err := sdk.AccAddressFromBech32(sNFT.SecondaryOwner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, fmt.Sprintf("bech32 address failed: %v", err))
	}
	if !msg.Sender.Equals(addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "nfts not associated with sender")
	}
	
	liveStream, found := keeper.GetLiveStream(ctx, msg.LiveStreamID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNFTNotFound, "liveStream not found")
	}
	
	programmeTimeDuration := 60 / liveStream.SlotsPerMinute
	startTime := msg.ProgramTime.Second()
	
	if (uint64(startTime) % programmeTimeDuration) != 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unable to create programme with %v duaration", msg.ProgramTime.String()))
	}
	
	dnftID := GenerateDNFTID(msg.ProgramTime)
	dnft, found := keeper.GetDNFT(ctx, dnftID)
	if !found {
		dnft.DNFTID = dnftID
		liveStream.DNFTIDs = append(liveStream.DNFTIDs, dnft.DNFTID)
	}
	
	if len(dnft.NFTID) < 1 && len(msg.SecondaryNFTID) > 0 {
		dnft.Type = TypeLocal
		dnft.NFTID = msg.SecondaryNFTID
		dnft.TwitterHandleName = sNFT.TwitterHandle
		dnft.TweetAssetID = sNFT.AssetID
	} else {
		return nil, sdkerrors.Wrap(types.ErrInvalidSlotBooking, "programme time already registered")
	}
	
	dnft.ProgramTime = msg.ProgramTime
	dnft.LiveStreamID = msg.LiveStreamID
	dnft.Status = "UNPAID"
	
	keeper.SetLiveStream(ctx, liveStream)
	keeper.SetDNFT(ctx, dnft)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMsgBookSlot,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeBookSlotValue),
			sdk.NewAttribute(types.AttributeSecondaryNFTID, dnft.NFTID),
			sdk.NewAttribute(types.AttributeDNFTID, dnft.DNFTID),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// Create LiveStream
func handleMsgLiveStreamCreation(ctx sdk.Context, keeper Keeper, msg MsgLiveStream) (*sdk.Result, error) {
	
	var liveStream BaseLiveStream
	
	account := keeper.GetAccount(ctx, msg.Sender)
	if account.GetAccountNumber() == 0 || account.GetAccountNumber() == 1 {
		count := keeper.GetGlobalLiveStreamCount(ctx)
		liveStreamID := GetLiveStreamID(count)
		liveStream, found := keeper.GetLiveStream(ctx, liveStreamID)
		if found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "liveStream already exits")
		}
		
		quo := 60 % msg.SlotsPerMin
		
		if quo != 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("unable to create livestream with these details %v", msg.SlotsPerMin))
		}
		
		liveStream.LiveStreamID = liveStreamID
		liveStream.CostPerAdPerSlot = msg.CostPerAdPerSlot
		liveStream.OwnerAddress = msg.Sender
		liveStream.Payout = msg.Payout
		liveStream.RevenueShare = msg.RevenueShare
		liveStream.SlotsPerMinute = msg.SlotsPerMin
		liveStream.Status = StatusActive
		
		keeper.SetLiveStream(ctx, liveStream)
		keeper.SetGlobalLiveStreamCount(ctx, count+1)
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "you are not authorised to create livestream")
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLiveStreamCreation,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeLiveStreamValue),
			sdk.NewAttribute(types.AttributeLiveStreamID, liveStream.LiveStreamID),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// Update liveStream

func handleMsgUpdateLiveStream(ctx sdk.Context, keeper Keeper, msg MsgUpdateLiveStream) (*sdk.Result, error) {
	
	liveStream, found := keeper.GetLiveStream(ctx, msg.LiveStreamID)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "liveStream not found")
	}
	
	if !liveStream.OwnerAddress.Equals(msg.Sender) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("livestream not associated with this %v address", msg.Sender))
	}
	
	liveStream.Payout = msg.Payout
	keeper.SetLiveStream(ctx, liveStream)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLiveStreamUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeLiveStreamUpdateValue),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
	
}

func handleMsgCreateAdNFT(ctx sdk.Context, keeper Keeper, msg MsgCreateAdNFT) (*sdk.Result, error) {
	count := keeper.GetGlobalAdsCount(ctx)
	adNFTId := GetAdNFTID(count)
	
	nft := BaseAdNFT{
		AssetID: msg.AssetID,
		Owner:   msg.Sender.String(),
		AdNFTID: adNFTId,
	}
	
	keeper.MintAdNFT(ctx, nft)
	keeper.SetAdNFTIDToAccount(ctx, msg.Sender, adNFTId)
	keeper.SetGlobalAdsCount(ctx, count+1)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvenTypeCreateAdNFT,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeAdNFTID, adNFTId),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
	
}

func handleMsgClaimTwitterAccount(ctx sdk.Context, k Keeper, msg MsgClaimTwitterAccount) (*sdk.Result, error) {
	
	info, _ := k.GetTwitterHandleInfo(ctx, msg.Handle)
	if reflect.DeepEqual(info, TwitterAccountInfo{}) {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFound, "account info doesn't exist")
	}
	
	if info.ClaimStatus {
		return nil, sdkerrors.Wrap(types.ErrInvalidClaimStatus, "account already claimed")
	}
	
	var count int64
	primaryNFTs := k.GetTweetsOfAccount(ctx, msg.PreOwner)
	for _, nft := range primaryNFTs {
		if strings.EqualFold(nft.TwitterHandle, msg.Handle) {
			count = +1
		}
	}
	
	if count == 0 {
		return nil, sdkerrors.Wrap(types.ErrAccountNotFound, "prev owner not associated with your twitter posts")
	}
	
	if err := k.SendCoins(ctx, msg.PreOwner, msg.Sender, info.LockedAmount); err != nil {
		return nil, err
	}
	
	info.ClaimStatus = true
	info.LockedAmount = sdk.Coins{}
	k.SetTwitterHandlerInfo(ctx, info)
	
	for _, nft := range primaryNFTs {
		if strings.EqualFold(nft.TwitterHandle, msg.Handle) {
			nft.PrimaryOwner = msg.Sender.String()
			k.UpdateTweetNFT(ctx, nft)
		}
		
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvenTypeClaimTwitterAccount,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.PreOwner.String()),
			sdk.NewAttribute(types.AttributeClaimantAddress, msg.Sender.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgUpdateAccessList(ctx sdk.Context, k Keeper, msg MsgUpdateAccessList) (*sdk.Result, error) {
	
	authorisedAddressed := k.GetAclAddressList(ctx)
	fromAccount := k.GetAccount(ctx, msg.Sender)
	if fromAccount.GetAccountNumber() == 0 || fromAccount.GetAccountNumber() == 1 {
		for _, addr := range authorisedAddressed.AccessList {
			if addr.Equals(msg.Address) {
				return nil, sdkerrors.Wrap(types.ErrAddressAlreadyExist, "")
			}
		}
		k.UpdateAclAddress(ctx, msg.Sender)
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "you are not authorised to update the access list")
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateAccessAddress,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeUpdatedAddress, msg.Address.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
	
}

func handleMsgUpdateHandlers(ctx sdk.Context, k Keeper, msg MsgUpdateHandlersInfo) (*sdk.Result, error) {
	
	authorisedHandlers := k.GetAuthorisedHandlerInfo(ctx)
	
	addressList := k.GetAclAddressList(ctx)
	for _, addr := range addressList.AccessList {
		if !addr.Equals(msg.Sender) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "you are not authorised to update handlers")
		}
	}
	
	for _, handle := range msg.Handlers {
		_, found := authorisedHandlers.Find(handle)
		if !found {
			
			authorisedHandlers.Handles = append(authorisedHandlers.Handles, handle)
		}
	}
	
	k.UpdaterHandlerInfo(ctx, authorisedHandlers)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateHandler,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
	
}
