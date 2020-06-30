package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
)

var (
	amino     = codec.New()
	ModuleCdc = codec.NewHybridCodec(amino, cdctypes.NewInterfaceRegistry())
)

// RegisterCodec registers the IBC transfer types
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgXNFTTransfer{}, "ibc/xnft/MsgXNFTTransfer", nil)
	cdc.RegisterConcrete(BaseNFTPacket{}, "ibc/xnft/BaseNFTPacket", nil)
	
	cdc.RegisterConcrete(MsgSlotBooking{}, "ibc/xnft/MsgSlotBooking", nil)
	cdc.RegisterConcrete(PacketSlotBooking{}, "ibc/xnft/PacketSlotBooking", nil)
	cdc.RegisterConcrete(PacketTokenDistribution{}, "ibc/xnft/PacketTokenDistribution", nil)
	cdc.RegisterConcrete(PacketPayLicensingFeeAndNFTTransfer{}, "ibc/xnft/PacketPayLicensingFeeAndNFTTransfer", nil)
	
	cdc.RegisterConcrete(MsgPayLicensingFee{}, "ibc/xnft/MsgPayLicensingFee", nil)
	cdc.RegisterConcrete(MsgSetParams{}, "xnft/MsgsetParams", nil)
	cdc.RegisterConcrete(Params{}, "xnft/Params", nil)
	cdc.RegisterConcrete(MsgDistributeFunds{}, "ibc/xnft/MsgDistributeFunds", nil)
	cdc.RegisterInterface((*XNFTs)(nil), nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgXNFTTransfer{})
	registry.RegisterImplementations((*XNFTs)(nil), &BaseNFTPacket{})
}

func init() {
	RegisterCodec(amino)
	channel.RegisterCodec(amino)
	commitmenttypes.RegisterCodec(amino)
	amino.Seal()
}
