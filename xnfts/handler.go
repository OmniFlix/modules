package xnfts

import (
	"strconv"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		
		switch msg := msg.(type) {
		case MsgXNFTTransfer:
			return handleMsgXNFTTransfer(ctx, k, msg)
		case MsgSlotBooking:
			return handleMsgSlotBooking(ctx, k, msg)
		case MsgSetParams:
			return handleMsgSetParams(ctx, k, msg)
		case MsgPayLicensingFee:
			return handlePayLicensingFeeAndNFTTransfer(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 xnft message type: %T", msg)
			
		}
	}
}

func handleMsgXNFTTransfer(ctx sdk.Context, k Keeper, msg MsgXNFTTransfer) (*sdk.Result, error) {
	
	err := k.XNFTTransfer(ctx, msg)
	if err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyReceiver, msg.Recipient),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgSlotBooking(ctx sdk.Context, k Keeper, msg MsgSlotBooking) (*sdk.Result, error) {
	var packet PacketSlotBooking
	
	if len(msg.AdNFTID) > 0 {
		_packet, err := k.BookSlotWithAdNFT(ctx, msg)
		if err != nil {
			return nil, err
		}
		packet = _packet
	}
	
	if len(msg.PrimaryNFTID) > 0 {
		_packet, err := k.BookSlotWithPrimaryNFT(ctx, msg)
		if err != nil {
			return nil, err
		}
		packet = _packet
	}
	
	packet.ProgramTime = msg.ProgrammeTime
	packet.LiveStreamID = msg.LiveStreamID
	
	if err := k.XTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, packet.GetBytes()); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleXNFTRecvPacket(ctx sdk.Context, k Keeper, packet channeltypes.Packet) (*sdk.Result, error) {
	
	var nftData BaseNFTPacket
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &nftData); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	
	acknowledgement := PostCreationPacketAcknowledgement{
		Success: true,
		Error:   "",
	}
	
	if err := k.OnRecvNFTPacket(ctx, nftData, packet); err != nil {
		acknowledgement = PostCreationPacketAcknowledgement{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeNFTPacketTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleSlotBookingRecvPacket(ctx sdk.Context, k Keeper, packet channeltypes.Packet) (*sdk.Result, error) {
	
	var slot PacketSlotBooking
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &slot); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	
	acknowledgement := PostCreationPacketAcknowledgement{
		Success: true,
		Error:   "",
	}
	
	if err := k.OnRecvSlotBooking(ctx, slot); err != nil {
		acknowledgement = PostCreationPacketAcknowledgement{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeNFTPacketTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleTokenDistributionRecvPacket(ctx sdk.Context, k Keeper, packet channeltypes.Packet) (*sdk.Result, error) {
	
	var data PacketTokenDistribution
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	
	acknowledgement := PostCreationPacketAcknowledgement{
		Success: true,
		Error:   "",
	}
	
	if err := k.OnRecvTokenDistribution(ctx, data); err != nil {
		acknowledgement = PostCreationPacketAcknowledgement{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeNFTPacketTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgSetParams(ctx sdk.Context, k Keeper, msg MsgSetParams) (*sdk.Result, error) {
	
	account := k.GetAccount(ctx, msg.Sender)
	if account.GetAccountNumber() == 0 || account.GetAccountNumber() == 1 {
		params := NewParams(msg.NFTChannel, msg.FFChannel, msg.DestHeight)
		k.SetParams(ctx, params)
		
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "")
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeSetParams,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyNFTChannel, msg.NFTChannel),
			sdk.NewAttribute(AttributeKeyDestHeight, strconv.Itoa(int(msg.DestHeight))),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handlePayLicensingFeeAndNFTTransfer(ctx sdk.Context, k Keeper, msg MsgPayLicensingFee) (*sdk.Result, error) {
	
	packet, err := k.PayLicensingFeeAndNFTTransfer(ctx, msg)
	if err != nil {
		return nil, err
	}
	if err := k.XTransfer(ctx, msg.SrcPort, msg.SrcChannel, msg.DestHeight, packet.GetBytes()); err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypePayLicensingFeeAndNFTTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Recipient),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.LicensingFee.String()),
		),
	)
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handlePayLicensingFeeAndNFTTransferRecvPacket(ctx sdk.Context, k Keeper, packet channeltypes.Packet) (*sdk.Result, error) {
	
	var data PacketPayLicensingFeeAndNFTTransfer
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	
	acknowledgement := PostCreationPacketAcknowledgement{
		Success: true,
		Error:   "",
	}
	
	if err := k.OnRecvXNFTTokenTransfer(ctx, data); err != nil {
		acknowledgement = PostCreationPacketAcknowledgement{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	nft, found := k.GetTweetNFTByID(ctx, data.PrimaryNFTID)
	if !found {
		acknowledgement = PostCreationPacketAcknowledgement{
			Success: false,
			Error:   "nfts not found",
		}
	}
	
	if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
		return nil, err
	}
	
	input := NFTInput{
		PrimaryNFTID:  data.PrimaryNFTID,
		Recipient:     data.Sender,
		AssetID:       nft.AssetID,
		LicensingFee:  data.LicensingFee,
		RevenueShare:  nft.RevenueShare,
		TwitterHandle: nft.TwitterHandle,
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypePayLicensingFeeAndNFTTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Recipient),
		),
	)
	
	msg := NewMsgXNFTTransfer(packet.DestinationPort, packet.DestinationChannel, packet.GetTimeoutHeight(),
		GetHexAddressFromBech32String(data.Recipient), input)
	
	if err := k.XNFTTransfer(ctx, msg); err != nil {
		return nil, err
	}
	
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
