package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, channelCap *capability.Capability, packet channelexported.PacketI) error
	PacketExecuted(ctx sdk.Context, chanCap *capability.Capability, packet channelexported.PacketI, acknowledgement []byte) error
	ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capability.Capability) error
}

type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capability.Capability
}

type (
	NFTKeeper interface {
		GetTweetNFTByID(ctx sdk.Context, id string) (nfts.BaseTweetNFT, bool)
		MintTweetNFT(ctx sdk.Context, nft nfts.BaseTweetNFT)
		GetAllTweetNFTs(ctx sdk.Context) []nfts.BaseTweetNFT
		GetTweetsOfAccount(ctx sdk.Context, address sdk.AccAddress) []nfts.BaseTweetNFT
		
		SetGlobalTweetCount(ctx sdk.Context, count uint64)
		GetGlobalTweetCount(ctx sdk.Context) uint64
		SetTweetIDToAccount(ctx sdk.Context, add sdk.AccAddress, id string)
	}
	
	BaseBankKeeper interface {
		AddCoins(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error)
		SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error)
	}
)
