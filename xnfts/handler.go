package xnfts

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	
	"github.com/FreeFlixMedia/modules/nfts"
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		
		switch msg := msg.(type) {
		case MsgXNFTTransfer:
			return handleMsgXNFTTransfer(ctx, k, msg)
		
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 xnft message type: %T", msg)
			
		}
	}
}

func handleMsgXNFTTransfer(ctx sdk.Context, k Keeper, msg MsgXNFTTransfer) (*sdk.Result, error) {
	var packet BaseNFTPacket
	
	if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext {
		nft, found := k.GetTweetNFTByID(ctx, msg.PrimaryNFTID)
		if !found {
			return nil, sdkerrors.Wrap(nfts.ErrNFTNotFound, "")
		}
		
		if !nft.License {
			return nil, sdkerrors.Wrap(nfts.ErrInvalidLicense, fmt.Sprintf("unable to transfer %s", nft.PrimaryNFTID))
		}
		
		if !msg.Sender.Equals(types.GetHexAddressFromBech32String(nft.PrimaryOwner)) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "")
		}
		
		packet.PrimaryNFTID = nft.PrimaryNFTID
		packet.PrimaryNFTOwner = nft.PrimaryOwner
		packet.License = nft.License
		packet.AssetID = nft.AssetID
		packet.RevenueShare = nft.RevenueShare
		packet.LicensingFee = nft.LicensingFee
		packet.SecondaryNFTOwner = msg.Recipient
		packet.TwitterHandle = nft.TwitterHandle
		
	} else if nfts.GetContextOfCurrentChain() == nfts.CoCoContext {
		
		count := k.GetGlobalTweetCount(ctx)
		sNFTID := nfts.GetSecondaryNFTID(count)
		
		packet.PrimaryNFTOwner = msg.Recipient
		packet.License = true
		packet.AssetID = msg.AssetID
		packet.RevenueShare = msg.RevenueShare
		packet.LicensingFee = msg.LicensingFee
		packet.SecondaryNFTID = sNFTID
		packet.SecondaryNFTOwner = msg.Sender.String()
		packet.TwitterHandle = msg.TwitterHandle
		
		k.MintTweetNFT(ctx, *packet.ToBaseTweetNFT())
		k.SetTweetIDToAccount(ctx, msg.Sender, sNFTID)
		k.SetGlobalTweetCount(ctx, count+1)
	}
	
	if err := k.XTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, packet.GetBytes()); err != nil {
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
