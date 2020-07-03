
### Codec

Once we have defined our messages, we need to describe to our encoder how they should be stored as bytes. To do this we edit the file located at `./nfts/internals/types/codec.go`. By describing our types as follows they will work with our encoding library:

```go=
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

```