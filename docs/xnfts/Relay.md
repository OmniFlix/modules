

### Relay.go

To relay the `NFTPacket`, we will access the channel created between two chains and the portId of the destination chain. Then we will create an outgoing packet that is `createOutgoingPacket`.

```go=
package keeper

import (
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
		addr, err := sdk.AccAddressFromBech32(data.PrimaryNFTOwner)
		if err != nil {
			return err
		}
		
		_, err = k.bankKeeper.AddCoins(ctx, addr, sdk.Coins{data.LicensingFee})
		if err != nil {
			return err
		}
		
		count := k.nftKeeper.GetGlobalTweetCount(ctx)
		primaryNFTID := nfts.GetPrimaryNFTID(count)
		data.PrimaryNFTID = primaryNFTID
		
		k.nftKeeper.MintTweetNFT(ctx, *data.ToBaseTweetNFT())
		k.SetTweetIDToAccount(ctx, addr, primaryNFTID)
		k.SetGlobalTweetCount(ctx, count+1)
		
		if err := k.XTransfer(ctx, packet.DestinationPort, packet.DestinationChannel, packet.TimeoutHeight, data.GetBytes()); err != nil {
			return err
		}
		
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
		
	}
	
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			nfts.EventTypeMsgMintTweetNFT,
			sdk.NewAttribute(sdk.AttributeKeySender, data.PrimaryNFTOwner),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.SecondaryNFTOwner),
			sdk.NewAttribute(nfts.AttributePrimaryNFTID, data.PrimaryNFTID),
			sdk.NewAttribute(nfts.AttributeSecondaryNFTID, data.SecondaryNFTID),
		),
	})
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

```
we have all of the basic actions of our module created, we want to make them accessible. We can do this with a CLI client and a REST client. For this tutorial we will just be creating a CLI client. If you are interested in what goes into making a REST client.

Let’s take a look at what goes into making a CLI.