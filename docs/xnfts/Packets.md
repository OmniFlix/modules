
### Packet Types

#### BaseNFTPacket

You may notice the reference to `types.BaseNFTPacket` throughout the `Keeper`. These is new struct defined in `./xnfts/types/packet.go` that contains all necessary information about different nft packets. You can create this file now and add the following:

```go=
package types

import (
	"encoding/json"
	"fmt"
	
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

type XNFTs interface {
	GetBytes() []byte
}

type BaseNFTPacket struct {
	PrimaryNFTID    string `json:"primary_nftid"`
	PrimaryNFTOwner string `json:"primary_nft_owner"`
	
	SecondaryNFTID    string `json:"secondary_nftid"`
	SecondaryNFTOwner string `json:"secondary_nft_owner"`
	
	AssetID string `json:"asset_id"`
	
	License      bool     `json:"license"`
	LicensingFee sdk.Coin `json:"licensing_fee"`
	RevenueShare sdk.Dec  `json:"revenue_share"`
	
	TwitterHandle string `json:"twitter_handle"`
}

var _ XNFTs = BaseNFTPacket{}

func (nft *BaseNFTPacket) Reset() {
	*nft = BaseNFTPacket{}
}

func (nft BaseNFTPacket) ProtoMessage() {
}

func NewBaseNFTPacket(primaryNFTID, secondaryNFTID, primaryNFTOwner, secondaryNFTOwner string,
	assetID, twitterHandle string, license bool, fee sdk.Coin, share sdk.Dec) BaseNFTPacket {
	return BaseNFTPacket{
		PrimaryNFTID:      primaryNFTID,
		PrimaryNFTOwner:   primaryNFTOwner,
		SecondaryNFTID:    secondaryNFTID,
		SecondaryNFTOwner: secondaryNFTOwner,
		AssetID:           assetID,
		License:           license,
		LicensingFee:      fee,
		RevenueShare:      share,
		TwitterHandle:     twitterHandle,
	}
}

func (nft BaseNFTPacket) String() string {
	return fmt.Sprintf(`
PrimaryNFTID: %s
PrimaryNFTOwner: %s

SecondaryNFTID: %s
SecondaryNFTOwner: %s

AssetID: %s
License: %t

LicensingFee: %s
RevenueShare: %s
TwittterHandle: %s
`, nft.PrimaryNFTID, nft.PrimaryNFTOwner, nft.SecondaryNFTID, nft.SecondaryNFTOwner, nft.AssetID, nft.License,
		nft.LicensingFee, nft.RevenueShare, nft.TwitterHandle)
}

func (nft BaseNFTPacket) ValidateBasic() error {
	if nft.PrimaryNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nftd should be present")
	} else if nft.PrimaryNFTOwner == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nft owner is empty")
	} else if nft.SecondaryNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nft id is empty")
	} else if nft.SecondaryNFTOwner == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nft owner is empty")
	} else if nft.AssetID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id is empty")
	} else if nft.RevenueShare.IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share is empty")
	} else if nft.TwitterHandle == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "handle name is empty")
	}
	return nil
}

func (nft BaseNFTPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(nft))
}

func (nft BaseNFTPacket) MarshalJSON() ([]byte, error) {
	type tmp BaseNFTPacket
	return json.Marshal(tmp(nft))
}

func (nft *BaseNFTPacket) UnmarshalJSON(bytes []byte) error {
	type tmp BaseNFTPacket
	var data tmp
	
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	
	*nft = BaseNFTPacket(data)
	return nil
}

func (nft BaseNFTPacket) ToBaseTweetNFT() *nfts.BaseTweetNFT {
	return &nfts.BaseTweetNFT{
		PrimaryNFTID:   nft.PrimaryNFTID,
		PrimaryOwner:   nft.PrimaryNFTOwner,
		SecondaryNFTID: nft.SecondaryNFTID,
		SecondaryOwner: nft.SecondaryNFTOwner,
		License:        nft.License,
		AssetID:        nft.AssetID,
		LicensingFee:   nft.LicensingFee,
		RevenueShare:   nft.RevenueShare,
		TwitterHandle:  nft.TwitterHandle,
	}
}

type PostCreationPacketAcknowledgement struct {
	Success bool   `json:"success" yaml:"success"`
	Error   string `json:"error" yaml:"error"`
}

func (ack PostCreationPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}

type PacketPayLicensingFeeAndNFTTransfer struct {
	PrimaryNFTID string   `json:"primary_nft_id"`
	LicensingFee sdk.Coin `json:"licensing_fee"`
	Recipient    string   `json:"recipient"`
	Sender       string   `json:"sender"`
}

func NewPacketPayLicensingFeeAndNFTTransfer(fee sdk.Coin, recipient, sender, primaryNFTID string) PacketPayLicensingFeeAndNFTTransfer {
	return PacketPayLicensingFeeAndNFTTransfer{
		LicensingFee: fee,
		PrimaryNFTID: primaryNFTID,
		Recipient:    recipient,
		Sender:       sender,
	}
}

var _ XNFTs = PacketPayLicensingFeeAndNFTTransfer{}

func (p PacketPayLicensingFeeAndNFTTransfer) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

func (p PacketPayLicensingFeeAndNFTTransfer) String() string {
	return fmt.Sprintf(`
PrimaryNFTID: %s,
Recipient: %s,
Sender: %s,
LicensingFee: %s
`, p.PrimaryNFTID, p.Recipient, p.Sender, p.LicensingFee.String())
}

func (p PacketPayLicensingFeeAndNFTTransfer) ValidateBasic() error {
	if len(p.PrimaryNFTID) == 0 {
		return fmt.Errorf("invalid input field, primary nfts id")
	}
	
	if p.LicensingFee.IsZero() {
		return fmt.Errorf("invalid licensing fee")
	}
	
	if len(p.Recipient) == 0 {
		return fmt.Errorf("invalid input field, recipient address")
	}
	if len(p.Sender) == 0 {
		return fmt.Errorf("invalid input field, sender address")
	}
	return nil
}

func (p PacketPayLicensingFeeAndNFTTransfer) MarshalJSON() ([]byte, error) {
	type tmp PacketPayLicensingFeeAndNFTTransfer
	return json.Marshal(tmp(p))
}

func (p *PacketPayLicensingFeeAndNFTTransfer) UnmarshalJSON(bytes []byte) error {
	type tmp PacketPayLicensingFeeAndNFTTransfer
	var data tmp
	
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	
	*p = PacketPayLicensingFeeAndNFTTransfer(data)
	return nil
}

```

You might also notice that each type has the `String` method. This allows us to render the struct as a string for rendering.

**ToBaseTweetNFT**

We use `ToBaseTweetNFT` to convert type from `BaseNFTPacket` to `BaseTweetNFT`
