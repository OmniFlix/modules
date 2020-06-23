package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMintTweetNFT{}, "nft/MsgMintTweetNFT", nil)
	cdc.RegisterConcrete(BaseTweetNFT{}, "nft/BaseTweetNFT", nil)
}

var (
	amino     = codec.New()
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
