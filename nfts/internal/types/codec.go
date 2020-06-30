package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateInitialTweetNFT{}, "nfts/MsgCreateInitialNFT", nil)
	cdc.RegisterConcrete(MsgMintTweetNFT{}, "nfts/MsgMintTweetNFT", nil)
	cdc.RegisterConcrete(MsgCreateAdNFT{}, "nfts/MsgCreateAdNFT", nil)
	cdc.RegisterConcrete(MsgLiveStream{}, "nfts/MsgLiveStream", nil)
	cdc.RegisterConcrete(MsgBookSlot{}, "nfts/MsgBookSlot", nil)
	cdc.RegisterConcrete(MsgUpdateAccessList{}, "nfts/MsgUpdateAccessList", nil)
	cdc.RegisterConcrete(MsgUpdateHandlersInfo{}, "nfts/MsgUpdateHandlersInfo", nil)
	cdc.RegisterConcrete(MsgClaimTwitterAccount{}, "nfts/MsgClaimTwitterAccount", nil)
	
	cdc.RegisterConcrete(BaseTweetNFT{}, "nfts/BaseTweetNFT", nil)
	cdc.RegisterConcrete(BaseAdNFT{}, "nfts/BaseAdNFT", nil)
	cdc.RegisterConcrete(BaseLiveStream{}, "nfts/BaseLiveStream", nil)
	cdc.RegisterConcrete(BaseDNFT{}, "nfts/BaseDNFT", nil)
}

var (
	amino     = codec.New()
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
