package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/golang/protobuf/proto"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

type NFTInput struct {
	PrimaryNFTID string `json:"primary_nftid"`
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
	
	// Packet BaseNFTPacket `json:"packet"`
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
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary nfts id is empty")
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

// -----------------------------------------------------------------

type MsgSlotBooking struct {
	PrimaryNFTID  string         `json:"primary_nftid"`
	AdNFTID       string         `json:"ad_nftid"`
	Sender        sdk.AccAddress `json:"sender"`
	Amount        sdk.Coin       `json:"amount"`
	ProgrammeTime time.Time      `json:"programme_time"`
	LiveStreamID  string         `json:"live_stream_id"`
	
	SourcePort    string `json:"source_port"`
	SourceChannel string `json:"source_channel"`
	DestHeight    uint64 `json:"dest_height"`
}

func NewMsgSlotBooking(primartNFTID, adNFTID, liveStreamID string, sender sdk.AccAddress, amount sdk.Coin, programmeTime time.Time, sourcePort, sourceChannel string, destHeight uint64) MsgSlotBooking {
	return MsgSlotBooking{
		PrimaryNFTID:  primartNFTID,
		AdNFTID:       adNFTID,
		LiveStreamID:  liveStreamID,
		Sender:        sender,
		Amount:        amount,
		ProgrammeTime: programmeTime,
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		DestHeight:    destHeight,
	}
}

var _ sdk.Msg = MsgSlotBooking{}

func (m *MsgSlotBooking) Reset() {
	*m = MsgSlotBooking{}
}

func (m *MsgSlotBooking) String() string {
	return proto.CompactTextString(m)
}

func (m MsgSlotBooking) ProtoMessage() {}

func (m MsgSlotBooking) Route() string {
	return RouterKey
}

func (m MsgSlotBooking) Type() string {
	return "msg_slot_booking"
}

func (m MsgSlotBooking) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(
		m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	
	if len(m.PrimaryNFTID) == 0 && len(m.AdNFTID) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "primary-nfts id or ad-nfts id  should not be empty")
	}
	if len(m.AdNFTID) > 1 {
		if !m.Amount.IsValid() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is invalid")
		}
	}
	
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	
	if len(m.LiveStreamID) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid livestream id")
	}
	if m.ProgrammeTime.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "program time should not be empty")
	}
	return nil
}

func (m MsgSlotBooking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSlotBooking) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

type MsgSetParams struct {
	Sender     sdk.AccAddress `json:"sender"`
	NFTChannel string         `json:"nft_channel"`
	FFChannel  string         `json:"ff_channel"`
	DestHeight uint64         `json:"dest_height"`
}

var _ sdk.Msg = MsgSetParams{}

func NewMsgSetParams(sender sdk.AccAddress, channel, ffchannel string, destHeight uint64) MsgSetParams {
	return MsgSetParams{
		Sender:     sender,
		NFTChannel: channel,
		FFChannel:  ffchannel,
		DestHeight: destHeight,
	}
}
func (m MsgSetParams) Route() string {
	return RouterKey
}

func (m MsgSetParams) Type() string {
	return "msg_set_params"
}

func (m MsgSetParams) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	
	if len(m.NFTChannel) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "nfts channel is empty")
	}
	
	if m.DestHeight < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "destination height not valid")
	}
	
	return nil
}

func (m MsgSetParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSetParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

type MsgPayLicensingFee struct {
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    string         `json:"recipient"`
	LicensingFee sdk.Coin       `json:"licensing_fee"`
	PrimaryNFTID string         `json:"primary_nftid"`
	
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

type MsgDistributeFunds struct {
	Channel    string         `json:"channel"`
	DestHeight uint64         `json:"dest_height"`
	Sender     sdk.AccAddress `json:"sender"`
}

func NewMsgDistributeFunds(channel string, height uint64, addr sdk.AccAddress) MsgDistributeFunds {
	return MsgDistributeFunds{
		Channel:    channel,
		DestHeight: height,
		Sender:     addr,
	}
}

var _ sdk.Msg = MsgDistributeFunds{}

func (m MsgDistributeFunds) Route() string {
	return RouterKey
}

func (m MsgDistributeFunds) Type() string {
	return "msg_distribute_funds"
}

func (m MsgDistributeFunds) ValidateBasic() error {
	if len(m.Channel) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "channel can't be empty")
	}
	
	if m.DestHeight < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "height can't be zero")
	}
	
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "address cant't be empty")
	}
	
	return nil
}

func (m MsgDistributeFunds) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDistributeFunds) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}
