
### Messages

Messages are a great place to start when building a module because they define the actions that your application can make. Think of all the scenarios where a user would be able to update the state of the application in any way. These should be boiled down into basic interactions, similar to **CRUD** (Create, Read, Update, Delete).

Let's start with **Create**

**MsgXNFTTransfer**

Messages are types which live inside the `./xnfts/internal/types/` directory. There is already a `msg.go` file. We can use this `msg.go`.

```go=
package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/golang/protobuf/proto"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

type NFTInput struct {
	PrimaryNFTID string `json:"primary_nft_id"`
	Recipient    string `json:"recipient"`
	AssetID      string `json:"asset_id"`
	
	LicensingFee  sdk.Coin `json:"licensing_fee"`
	RevenueShare  sdk.Dec  `json:"revenue_share"`
	TwitterHandle string   `json:"twitter_handle"`
}

type MsgXNFTTransfer struct {
	SourcePort    string         `json:"source_port"`
	SourceChannel string         `json:"source_channel"`
	DestHeight    uint64         `json:"dest_height"`
	Sender        sdk.AccAddress `json:"sender"`
	
	NFTInput
}

func NewMsgXNFTTransfer(sourcePort, sourceChannel string, height uint64, sender sdk.AccAddress,
	nftInput NFTInput) MsgXNFTTransfer {
	return MsgXNFTTransfer{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		DestHeight:    height,
		Sender:        sender,
		NFTInput:      nftInput,
	}
}

var _ sdk.Msg = MsgXNFTTransfer{}

func (m *MsgXNFTTransfer) Reset() {
	*m = MsgXNFTTransfer{}
}

func (m *MsgXNFTTransfer) String() string {
	return proto.CompactTextString(m)
}

func (m MsgXNFTTransfer) ProtoMessage() {}

func (m MsgXNFTTransfer) Route() string {
	return RouterKey
}

func (m MsgXNFTTransfer) Type() string {
	return "msg_xnft_transfer"
}

func (m MsgXNFTTransfer) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	
	if nfts.GetContextOfCurrentChain() == nfts.CoCoContext {
		if m.NFTInput.AssetID == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
		} else if m.NFTInput.RevenueShare.IsZero() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share is not allowed to be empty")
		} else if m.NFTInput.TwitterHandle == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "handle name should not be empty")
		} else if !m.NFTInput.LicensingFee.IsValid() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "licensing fee is invalid")
		}
	}
	
	if nfts.GetContextOfCurrentChain() == nfts.FreeFlixContext {
		if len(m.PrimaryNFTID) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nft id is empty")
		}
	}
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	} else if m.NFTInput.Recipient == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Recipient should not be nil")
	}
	return nil
}

func (m MsgXNFTTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgXNFTTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

// --------------------------------------------------------------------

type MsgPayLicensingFee struct {
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    string         `json:"recipient"`
	LicensingFee sdk.Coin       `json:"licensing_fee"`
	PrimaryNFTID string         `json:"primary_nft_id"`
	
	SrcPort    string `json:"src_port"`
	SrcChannel string `json:"src_channel"`
	DestHeight uint64 `json:"dest_height"`
}

func NewMsgPayLicensingFee(
	sourcePort, sourceChannel, primaryNFTID string, destHeight uint64, fee sdk.Coin, sender sdk.AccAddress, receiver string,
) MsgPayLicensingFee {
	return MsgPayLicensingFee{
		SrcPort:      sourcePort,
		SrcChannel:   sourceChannel,
		DestHeight:   destHeight,
		PrimaryNFTID: primaryNFTID,
		LicensingFee: fee,
		Sender:       sender,
		Recipient:    receiver,
	}
}

var _ sdk.Msg = MsgPayLicensingFee{}

func (m MsgPayLicensingFee) Route() string {
	return RouterKey
}

func (m MsgPayLicensingFee) Type() string {
	return "msg_pay_licensing_fee_and_nft_transfer"
}

func (m MsgPayLicensingFee) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SrcPort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SrcChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	
	if !m.LicensingFee.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if m.Recipient == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	
	if len(m.PrimaryNFTID) == 0 {
		return fmt.Errorf("invalid input field, primary nfts id")
	}
	
	if m.LicensingFee.IsZero() {
		return fmt.Errorf("invalid licensing fee")
	}
	
	return nil
}

func (m MsgPayLicensingFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgPayLicensingFee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}


```

Notice that all Messages in the app need to follow the `sdk.Msg` interface. The Message `struct` contains all the necessary information when creating a new xnft:

* `Sender`- Who initiates xnft transfer. This uses the `sdk.AccAddress` type which represents an account in the app.
* `SourcePort`- source chain port
* `SourceChannel`- source chain channel
* `DestHeight`- Destination chain height
* `NFTInput`:
    * `PrimaryNFTID`- Primary TweetNFT ID
    * `Recipient`- Primary nft receiver
    * `AssetID`- AssetID of the tweet nft.
    * `LicensingFee`- Amount to pay to get the license of an asset
    * `RevenueShare`- Amount of share we get when nft is used
    * `TwitterHandle`- Twitter handle of the user.


The `Msg` interface requires some other methods to be set, like validating the content of the `struct`, and confirming the msg was signed and submitted by the Creator.
