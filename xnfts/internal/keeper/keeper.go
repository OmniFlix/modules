package keeper

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/tendermint/tendermint/libs/log"
	
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
