package keeper

import (
	"fmt"
	"strings"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	
	"github.com/FreeFlixMedia/modules/nfts"
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

const (
	DefaultPacketTimeoutHeight = 1000
	
	DefaultPacketTimeoutTimestamp = 0
)

func (k Keeper) XTransfer(
	ctx sdk.Context,
	sourcePort, sourceChannel string,
	destHeight uint64,
	packetData []byte,
) error {
	
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}
	
	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()
	
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channeltypes.ErrSequenceSendNotFound
	}
	
	return k.createOutgoingPacket(ctx, sequence, sourcePort, sourceChannel, destinationPort, destinationChannel, destHeight, packetData)
}

func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	sourcePort, sourceChannel string,
	destinationPort, destinationChannel string,
	destHeight uint64,
	data []byte,
) error {
	
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}
	
	packet := channeltypes.NewPacket(
		data,
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeoutHeight, // TODO : DestHeight need to be updated with src header.height
		DefaultPacketTimeoutTimestamp,
	)
	
	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) OnRecvNFTPacket(ctx sdk.Context, data types.BaseNFTPacket, packet channeltypes.Packet) error {
	
	if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext && len(data.PrimaryNFTID) == 0 {
		
		authorisedHandlers := k.nftKeeper.GetAuthorisedHandlerInfo(ctx)
		_, found := authorisedHandlers.Find(data.TwitterHandle)
		if !found {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised handler to create asset")
		}
		
		addr, err := sdk.AccAddressFromBech32(data.PrimaryNFTOwner)
		if err != nil {
			return err
		}
		
		getTwitterAccount, found := k.GetTwitterHandleInfo(ctx, data.TwitterHandle)
		if found {
			_, err = k.bankKeeper.AddCoins(ctx, getTwitterAccount.Owner, sdk.Coins{data.LicensingFee})
			if err != nil {
				return err
			}
			
			tweetNFTs := k.nftKeeper.GetTweetsOfAccount(ctx, getTwitterAccount.Owner)
			
			for _, nt := range tweetNFTs {
				if strings.EqualFold(nt.AssetID, data.AssetID) {
					return sdkerrors.Wrap(nfts.ErrAssetIDAlreadyExist, "")
				}
			}
			
		} else {
			prevAmount := getTwitterAccount.LockedAmount
			getTwitterAccount.LockedAmount = prevAmount.Add(data.LicensingFee)
			getTwitterAccount.Handle = data.TwitterHandle
			
			// Add coins to FF Default Address
			_, err = k.bankKeeper.AddCoins(ctx, addr, sdk.Coins{data.LicensingFee})
			if err != nil {
				return err
			}
		}
		
		k.SetTwitterHandlerInfo(ctx, getTwitterAccount)
		count := k.nftKeeper.GetGlobalTweetCount(ctx)
		primaryNFTID := nfts.GetPrimaryNFTID(count)
		data.PrimaryNFTID = primaryNFTID
		
		k.nftKeeper.MintTweetNFT(ctx, *data.ToBaseTweetNFT())
		k.SetTweetIDToAccount(ctx, addr, primaryNFTID)
		k.SetGlobalTweetCount(ctx, count+1)
		
		if err := k.XTransfer(ctx, packet.DestinationPort, packet.DestinationChannel, packet.TimeoutHeight, data.GetBytes()); err != nil {
			return err
		}
		
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeMintNFT,
				sdk.NewAttribute(sdk.AttributeKeySender, data.SecondaryNFTOwner),
				sdk.NewAttribute(types.AttributeKeyReceiver, data.PrimaryNFTOwner),
				sdk.NewAttribute(types.AtttibutePrimaryNFTID, data.PrimaryNFTID),
				sdk.NewAttribute(types.AtttibuteSecondaryNFTID, data.SecondaryNFTID),
			),
		})
		
	}
	if nfts.GetContextOfCurrentChain() == nfts.CoCoContext && len(data.SecondaryNFTID) == 0 {
		addr, err := sdk.AccAddressFromBech32(data.SecondaryNFTOwner)
		if err != nil {
			return err
		}
		
		count := k.nftKeeper.GetGlobalTweetCount(ctx)
		secondaryNFTID := nfts.GetSecondaryNFTID(count)
		data.SecondaryNFTID = secondaryNFTID
		
		k.nftKeeper.MintTweetNFT(ctx, *data.ToBaseTweetNFT())
		k.nftKeeper.SetTweetIDToAccount(ctx, addr, secondaryNFTID)
		k.nftKeeper.SetGlobalTweetCount(ctx, count+1)
		
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeMintNFT,
				sdk.NewAttribute(sdk.AttributeKeySender, data.PrimaryNFTOwner),
				sdk.NewAttribute(types.AttributeKeyReceiver, data.SecondaryNFTOwner),
				sdk.NewAttribute(types.AtttibutePrimaryNFTID, data.PrimaryNFTID),
				sdk.NewAttribute(types.AtttibuteSecondaryNFTID, data.SecondaryNFTID),
			),
		})
		
	}
	
	if len(data.PrimaryNFTID) >= 0 && len(data.SecondaryNFTID) >= 0 {
		k.nftKeeper.UpdateTweetNFT(ctx, *data.ToBaseTweetNFT())
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateNFT,
				sdk.NewAttribute(sdk.AttributeKeySender, data.PrimaryNFTOwner),
				sdk.NewAttribute(types.AttributeKeyReceiver, data.SecondaryNFTOwner),
				sdk.NewAttribute(types.AtttibutePrimaryNFTID, data.PrimaryNFTID),
				sdk.NewAttribute(types.AtttibuteSecondaryNFTID, data.SecondaryNFTID),
			),
		})
	}
	
	return nil
}

// -------------------------------------------

func (k Keeper) OnRecvSlotBooking(ctx sdk.Context, data types.PacketSlotBooking) error {
	
	liveStream, ok := k.nftKeeper.GetLiveStream(ctx, data.LiveStreamID)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "liveStream not found.")
	}
	
	programmeTimeDuration := 60 / liveStream.SlotsPerMinute
	startTime := data.ProgramTime.Second()
	
	if (uint64(startTime) % programmeTimeDuration) != 0 {
		return fmt.Errorf(" invalid startime for programme %s", data.ProgramTime)
	}
	
	dnftID := nfts.GenerateDNFTID(data.ProgramTime)
	dnft, found := k.nftKeeper.GetDNFT(ctx, dnftID)
	if !found {
		dnft.DNFTID = dnftID
		dnft.Status = "UNPAID"
		liveStream.DNFTIDs = append(liveStream.DNFTIDs, dnftID)
	}
	
	if data.AdNFTID != "" && dnft.AdNFTID == "" {
		dnft.AdNFTID = data.AdNFTID
		dnft.LockedAmount = data.Amount
		dnft.AdNFTAssetID = data.AdNFTAssetID
	}
	
	if data.PrimaryNFTID != "" && dnft.NFTID == "" {
		dnft.NFTID = data.PrimaryNFTID
		dnft.PrimaryNFTAddress = data.PrimaryNFTOwner
		dnft.TweetAssetID = data.TweetNFTAssetID
		dnft.Type = nfts.TypeIBC
		dnft.TwitterHandleName = data.HandleName
	}
	
	dnft.ProgramTime = data.ProgramTime
	dnft.LiveStreamID = data.LiveStreamID
	k.nftKeeper.SetDNFT(ctx, dnft)
	k.nftKeeper.SetLiveStream(ctx, liveStream)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			nfts.EventTypeMsgBookSlot,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(nfts.AttributeDNFTID, dnft.DNFTID),
		),
	)
	return nil
}

func (k Keeper) OnRecvTokenDistribution(ctx sdk.Context, data types.PacketTokenDistribution) error {
	
	receiver, _ := sdk.AccAddressFromBech32(data.Recipient)
	
	getTweetAccount, found := k.nftKeeper.GetTwitterHandleInfo(ctx, data.Handler)
	if !found {
		prevAmount := getTweetAccount.LockedAmount
		getTweetAccount.LockedAmount = prevAmount.Add(data.AmountLocked)
		getTweetAccount.ClaimStatus = false
		getTweetAccount.Handle = data.Handler
		_, err := k.bankKeeper.AddCoins(ctx, receiver, sdk.Coins{data.AmountLocked})
		if err != nil {
			return err
		}
	} else {
		_, err := k.bankKeeper.AddCoins(ctx, getTweetAccount.Owner, sdk.Coins{data.AmountLocked})
		if err != nil {
			return err
		}
	}
	
	k.nftKeeper.SetTwitterHandlerInfo(ctx, getTweetAccount)
	
	return nil
}

func (k Keeper) OnRecvXNFTTokenTransfer(ctx sdk.Context, data types.PacketPayLicensingFeeAndNFTTransfer) error {
	
	receiver, err := sdk.AccAddressFromBech32(data.Recipient)
	if err != nil {
		return err
	}
	
	_, err = k.bankKeeper.AddCoins(ctx, receiver, sdk.Coins{data.LicensingFee})
	if err != nil {
		return err
	}
	return nil
}
