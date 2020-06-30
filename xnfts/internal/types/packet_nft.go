package types

import (
	"encoding/json"
	"fmt"
	"time"
	
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	nnft "github.com/FreeFlixMedia/modules/nfts"
)

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

func (nft *BaseNFTPacket) Reset() {
	*nft = BaseNFTPacket{}
}

func (nft BaseNFTPacket) ProtoMessage() {
}

var _ XNFTs = BaseNFTPacket{}

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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nfts owner is empty")
	} else if nft.SecondaryNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nfts id is empty")
	} else if nft.SecondaryNFTOwner == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nfts owner is empty")
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

func (nft BaseNFTPacket) ToBaseTweetNFT() *nnft.BaseTweetNFT {
	return &nnft.BaseTweetNFT{
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

type PacketSlotBooking struct {
	PrimaryNFTID    string    `json:"primary_nftid"`
	PrimaryNFTOwner string    `json:"primary_nft_owner"`
	TweetNFTAssetID string    `json:"tweet_nft_asset_id"`
	HandleName      string    `json:"handle_name"`
	AdNFTID         string    `json:"ad_nftid"`
	AdNFTAssetID    string    `json:"ad_nft_asset_id"`
	ProgramTime     time.Time `json:"program_time"`
	Amount          sdk.Coin  `json:"amount"`
	LiveStreamID    string    `json:"live_stream_id"`
}

func NewPacketSlotBooking(primaryNFTID, adNFTID, liveStreamID, primaryNFTOwner, tweetNFTAssetID, adNFTAssetID, handle string, programTime time.Time, amount sdk.Coin) PacketSlotBooking {
	return PacketSlotBooking{
		PrimaryNFTID:    primaryNFTID,
		PrimaryNFTOwner: primaryNFTOwner,
		HandleName:      handle,
		TweetNFTAssetID: tweetNFTAssetID,
		AdNFTAssetID:    adNFTAssetID,
		AdNFTID:         adNFTID,
		ProgramTime:     programTime,
		LiveStreamID:    liveStreamID,
		Amount:          amount,
	}
}

var _ XNFTs = PacketSlotBooking{}

func (p PacketSlotBooking) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

func (p PacketSlotBooking) String() string {
	return fmt.Sprintf(`
PrimaryNFTID: %s,
AdNFTID: %s,
ProgramTime: %s
`, p.PrimaryNFTID, p.AdNFTID, p.ProgramTime)
}

func (p PacketSlotBooking) ValidateBasic() error {
	if p.AdNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "adnft id should not be empty")
	} else if p.PrimaryNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nfts id should not be empty")
	} else if p.ProgramTime.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "program time should not be empty")
	} else if !p.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is invalid")
	} else if p.LiveStreamID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "live streamID should not be empty")
	}
	return nil
}

func (p PacketSlotBooking) MarshalJSON() ([]byte, error) {
	type tmp PacketSlotBooking
	return json.Marshal(tmp(p))
}

func (p *PacketSlotBooking) UnmarshalJSON(bytes []byte) error {
	type tmp PacketSlotBooking
	var data tmp
	
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	
	*p = PacketSlotBooking(data)
	return nil
}

type PacketTokenDistribution struct {
	AmountLocked sdk.Coin `json:"amount_locked"`
	Recipient    string   `json:"recipient"`
	Handler      string   `json:"handler"`
	Sender       string   `json:"sender"`
}

func NewPacketTokenTransfer(amount sdk.Coin, recipient, sender, handler string) PacketTokenDistribution {
	return PacketTokenDistribution{
		AmountLocked: amount,
		Recipient:    recipient,
		Sender:       sender,
		Handler:      handler,
	}
}

var _ XNFTs = PacketTokenDistribution{}

func (p PacketTokenDistribution) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

func (p PacketTokenDistribution) String() string {
	return fmt.Sprintf(`
AmountLocked: %s,
Recipient: %s,
Sender: %s,
Handler: %s,
`, p.AmountLocked, p.Recipient, p.Sender, p.Handler)
}

func (p PacketTokenDistribution) ValidateBasic() error {
	return nil // TODO
}

func (p PacketTokenDistribution) MarshalJSON() ([]byte, error) {
	type tmp PacketTokenDistribution
	return json.Marshal(tmp(p))
}

func (p *PacketTokenDistribution) UnmarshalJSON(bytes []byte) error {
	type tmp PacketTokenDistribution
	var data tmp
	
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	
	*p = PacketTokenDistribution(data)
	return nil
}

type PacketPayLicensingFeeAndNFTTransfer struct {
	PrimaryNFTID string   `json:"primary_nftid"`
	LicensingFee sdk.Coin `json:"licensing_fee"`
	Recipient    string   `json:"recipient"`
	Sender       string   `json:"sender"`
}

func NewPacketXNFTTokenTransfer(fee sdk.Coin, recipient, sender, primaryNFTID string) PacketPayLicensingFeeAndNFTTransfer {
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
