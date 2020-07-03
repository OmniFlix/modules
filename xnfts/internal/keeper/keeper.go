package keeper

import (
	"fmt"
	"strings"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/FreeFlixMedia/modules/nfts"
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	nftKeeper types.NFTKeeper
	
	bankKeeper    types.BaseBankKeeper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	scopedKeeper  capability.ScopedKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, nftKeeper types.NFTKeeper, bankKeeper types.BaseBankKeeper, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	scopedKeeper capability.ScopedKeeper) Keeper {
	
	return Keeper{
		
		storeKey:      key,
		cdc:           cdc,
		nftKeeper:     nftKeeper,
		bankKeeper:    bankKeeper,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.ModuleName))
}

func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
	}
	
	return k.channelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(portID, channelID))
	if !ok {
		return sdkerrors.Wrapf(channel.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", ibctypes.ChannelCapabilityPath(portID, channelID))
	}
	
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

func (k Keeper) IsBounded(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.PortPath(portID))
	return ok
}

func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	store := ctx.KVStore(k.storeKey)
	
	store.Set([]byte(types.PortKey), []byte(portID))
	capa := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capa, ibctypes.PortPath(portID))
}

func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (keeper Keeper) PayLicensingFeeAndNFTTransfer(ctx sdk.Context, msg types.MsgPayLicensingFee) (
	types.PacketPayLicensingFeeAndNFTTransfer, error) {
	snfts := keeper.GetAllTweetNFTs(ctx)
	
	for _, _nft := range snfts {
		if strings.EqualFold(_nft.PrimaryNFTID, msg.PrimaryNFTID) {
			return types.PacketPayLicensingFeeAndNFTTransfer{}, sdkerrors.Wrap(nfts.ErrInvalidLicense, "primary nfts already licensed")
		}
	}
	
	_, err := keeper.SubtractCoins(ctx, msg.Sender, sdk.Coins{msg.LicensingFee})
	if err != nil {
		return types.PacketPayLicensingFeeAndNFTTransfer{}, err
	}
	
	packet := types.PacketPayLicensingFeeAndNFTTransfer{
		PrimaryNFTID: msg.PrimaryNFTID,
		LicensingFee: msg.LicensingFee,
		Sender:       msg.Sender.String(),
		Recipient:    msg.Recipient,
	}
	
	return packet, nil
}

func (keeper Keeper) UpdateSecondaryNFTOwner(ctx sdk.Context, msg types.MsgXNFTTransfer) (types.BaseNFTPacket, error) {
	var packet types.BaseNFTPacket
	_nft, found := keeper.GetTweetNFTByID(ctx, msg.PrimaryNFTID)
	if !found {
		return types.BaseNFTPacket{}, sdkerrors.Wrap(nfts.ErrNFTNotFound, "")
	}
	
	if !_nft.License {
		return types.BaseNFTPacket{}, sdkerrors.Wrap(nfts.ErrInvalidLicense, fmt.Sprintf("unable to transfer %s", _nft.PrimaryNFTID))
	}
	
	if !msg.Sender.Equals(types.GetHexAddressFromBech32String(_nft.PrimaryOwner)) {
		return types.BaseNFTPacket{}, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "")
	}
	
	packet.PrimaryNFTID = _nft.PrimaryNFTID
	packet.PrimaryNFTOwner = _nft.PrimaryOwner
	packet.License = _nft.License
	packet.AssetID = _nft.AssetID
	packet.RevenueShare = _nft.RevenueShare
	packet.LicensingFee = _nft.LicensingFee
	packet.SecondaryNFTOwner = msg.Recipient
	packet.TwitterHandle = _nft.TwitterHandle
	
	return packet, nil
}

func (keeper Keeper) CreateSecondaryNFT(ctx sdk.Context, msg types.MsgXNFTTransfer) (types.BaseNFTPacket, error) {
	var packet types.BaseNFTPacket
	
	_, err := keeper.SubtractCoins(ctx, msg.Sender, sdk.Coins{msg.LicensingFee})
	if err != nil {
		return types.BaseNFTPacket{}, err
	}
	
	count := keeper.GetGlobalTweetCount(ctx)
	sNFTID := nfts.GetSecondaryNFTID(count)
	
	packet.PrimaryNFTOwner = msg.Recipient
	packet.License = true
	packet.AssetID = msg.AssetID
	packet.RevenueShare = msg.RevenueShare
	packet.LicensingFee = msg.LicensingFee
	packet.SecondaryNFTID = sNFTID
	packet.SecondaryNFTOwner = msg.Sender.String()
	packet.TwitterHandle = msg.TwitterHandle
	
	keeper.MintTweetNFT(ctx, *packet.ToBaseTweetNFT())
	keeper.SetTweetIDToAccount(ctx, msg.Sender, sNFTID)
	keeper.SetGlobalTweetCount(ctx, count+1)
	
	return packet, nil
}

func (keeper Keeper) XNFTTransfer(ctx sdk.Context, msg types.MsgXNFTTransfer) error {
	var packet types.BaseNFTPacket
	
	if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext {
		_packet, err := keeper.UpdateSecondaryNFTOwner(ctx, msg)
		if err != nil {
			return err
		}
		packet = _packet
		
	} else if nfts.GetContextOfCurrentChain() == nfts.CoCoContext {
		_packet, err := keeper.CreateSecondaryNFT(ctx, msg)
		if err != nil {
			return err
		}
		packet = _packet
	}
	
	if err := keeper.XTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, packet.GetBytes()); err != nil {
		return err
	}
	
	return nil
}
