package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MsgTypeCreateInitialTweetNFT = "msg_create_initial_tweet_nft"
	MsgTypeMintNFT               = "msg_mint_nft"
	MsgTypeCreateAdNFT           = "msg_ad_nft"
	MsgTypeClaimTwitterAccount   = "msg_claim_twitter_account"
)

type MsgCreateInitialTweetNFT struct {
	Sender        sdk.AccAddress `json:"sender"`
	AssetID       string         `json:"asset_id"`
	TwitterHandle string         `json:"twitter_handle"`
}

var _ sdk.Msg = MsgCreateInitialTweetNFT{}

func NewMsgCreateInitialTweetNFT(sender sdk.AccAddress, assetID, handle string) MsgCreateInitialTweetNFT {
	return MsgCreateInitialTweetNFT{
		Sender:        sender,
		AssetID:       assetID,
		TwitterHandle: handle,
	}
}

func (m MsgCreateInitialTweetNFT) Route() string {
	return RouterKey
}

func (m MsgCreateInitialTweetNFT) Type() string {
	return MsgTypeCreateInitialTweetNFT
}

func (m MsgCreateInitialTweetNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	} else if m.AssetID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
	} else if m.TwitterHandle == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "twitter handle should not be empty")
	}
	return nil
}

func (m MsgCreateInitialTweetNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateInitialTweetNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

// -----------------------------------------------------------------
type MsgMintTweetNFT struct {
	Sender        sdk.AccAddress `json:"sender"`
	AssetID       string         `json:"asset_id"`
	License       bool           `json:"license"`
	LicensingFee  sdk.Coin       `json:"licensing_fee"`
	RevenueShare  sdk.Dec        `json:"revenue_share"`
	TwitterHandle string         `json:"twitter_handle"`
}

func NewMsgMintNFT(sender sdk.AccAddress, assetID string, license bool, fee sdk.Coin, share sdk.Dec, handle string) MsgMintTweetNFT {
	return MsgMintTweetNFT{
		Sender:        sender,
		AssetID:       assetID,
		License:       license,
		LicensingFee:  fee,
		RevenueShare:  share,
		TwitterHandle: handle,
	}
}

var _ sdk.Msg = MsgMintTweetNFT{}

func (m MsgMintTweetNFT) Route() string {
	return RouterKey
}

func (m MsgMintTweetNFT) Type() string {
	return MsgTypeMintNFT
}

func (m MsgMintTweetNFT) ValidateBasic() error {
	
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	} else if m.AssetID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
	}
	
	if m.License {
		if m.LicensingFee.IsZero() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid licensing fee provided")
		} else if m.RevenueShare.IsZero() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "revenue share should not be nil")
		}
	}
	if m.TwitterHandle == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "twitter handle should not be empty")
	}
	return nil
}

func (m MsgMintTweetNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgMintTweetNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

// ----------------------------------------------------------------------
// MsgLiveStream

type MsgLiveStream struct {
	Sender           sdk.AccAddress `json:"sender"`
	RevenueShare     sdk.Dec        `json:"revenue_share"`
	CostPerAdPerSlot sdk.Coin       `json:"cost_per_ad_per_slot"`
	Payout           time.Duration  `json:"payout"`
	SlotsPerMin      uint64         `json:"slots_per_min"`
}

func NewMsgLiveStream(sender sdk.AccAddress, share sdk.Dec, costPerSlot sdk.Coin,
	payout time.Duration, slotsPerMin uint64) MsgLiveStream {
	return MsgLiveStream{
		Sender:           sender,
		RevenueShare:     share,
		CostPerAdPerSlot: costPerSlot,
		Payout:           payout,
		SlotsPerMin:      slotsPerMin,
	}
}

var _ sdk.Msg = MsgLiveStream{}

func (msg MsgLiveStream) Route() string { return RouterKey }

func (msg MsgLiveStream) Type() string { return "msg_create_live_stream" }

func (msg MsgLiveStream) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	} else if msg.RevenueShare.IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid revenue share provided")
	} else if msg.Payout == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "payout time should more than zero")
	} else if !msg.CostPerAdPerSlot.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "cost per ad per slot is invalid")
	} else if msg.SlotsPerMin == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "slots per min should be more than zero")
	}
	return nil
}

func (msg MsgLiveStream) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgLiveStream) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// --------------------------------------------------------------------------
// update liveStream msg

// Update liveStream

type MsgUpdateLiveStream struct {
	Sender       sdk.AccAddress `json:"sender"`
	Payout       time.Duration  `json:"payout"`
	LiveStreamID string         `json:"live_stream_id"`
}

func NewUpdateLiveStream(addr sdk.AccAddress, payout time.Duration, id string) MsgUpdateLiveStream {
	return MsgUpdateLiveStream{
		Sender:       addr,
		Payout:       payout,
		LiveStreamID: id,
	}
}

var _ sdk.Msg = MsgUpdateLiveStream{}

func (msg MsgUpdateLiveStream) Route() string { return RouterKey }

func (msg MsgUpdateLiveStream) Type() string { return "msg_update_live_stream" }

func (msg MsgUpdateLiveStream) ValidateBasic() error {
	if msg.Payout == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "payout time should be more than zero")
	} else if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address is invalid")
	} else if msg.LiveStreamID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "live stream should not be empty")
	}
	return nil
}

func (msg MsgUpdateLiveStream) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUpdateLiveStream) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -----------------------------------------------------------------------
// dnft creation msg

type MsgBookSlot struct {
	Sender         sdk.AccAddress `json:"sender"`
	SecondaryNFTID string         `json:"secondary_nftid"`
	ProgramTime    time.Time      `json:"program_time"`
	LiveStreamID   string         `json:"live_stream_id"`
}

func NewMsgBookSlot(sender sdk.AccAddress, sNFTID, streamID string, programTime time.Time) MsgBookSlot {
	return MsgBookSlot{
		Sender:         sender,
		SecondaryNFTID: sNFTID,
		ProgramTime:    programTime,
		LiveStreamID:   streamID,
	}
}

var _ sdk.Msg = MsgBookSlot{}

func (msg MsgBookSlot) Route() string { return RouterKey }
func (msg MsgBookSlot) Type() string {
	return "msg_slot_booking"
}

func (msg MsgBookSlot) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	} else if msg.LiveStreamID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "live stream id should not be empty")
	} else if msg.ProgramTime.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "program time should not be zero")
	} else if msg.SecondaryNFTID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secondary nfts id shoud not be empty")
	}
	return nil
}

func (msg MsgBookSlot) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBookSlot) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -------------------------------------------------

type MsgCreateAdNFT struct {
	AssetID string         `json:"asset_id"`
	Sender  sdk.AccAddress `json:"sender"`
}

func NewMsgCreateAdNFT(assetID string, sender sdk.AccAddress) MsgCreateAdNFT {
	return MsgCreateAdNFT{
		AssetID: assetID,
		Sender:  sender,
	}
}

var _ sdk.Msg = MsgCreateAdNFT{}

func (m MsgCreateAdNFT) Route() string {
	return RouterKey
}

func (m MsgCreateAdNFT) Type() string {
	return MsgTypeCreateAdNFT
}

func (m MsgCreateAdNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender addr")
	} else if m.AssetID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset id should not be empty")
	}
	return nil
}

func (m MsgCreateAdNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateAdNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

// ----------------------------------------

type MsgClaimTwitterAccount struct {
	Sender   sdk.AccAddress `json:"sender"`
	Handle   string         `json:"handle"`
	PreOwner sdk.AccAddress `json:"pre_owner"`
}

var _ sdk.Msg = MsgClaimTwitterAccount{}

func NewMsgClaimTwitterAccount(sender, prevOwner sdk.AccAddress, handle string) MsgClaimTwitterAccount {
	return MsgClaimTwitterAccount{
		Sender:   sender,
		Handle:   handle,
		PreOwner: prevOwner,
	}
}

func (m MsgClaimTwitterAccount) Route() string {
	return RouterKey
}

func (m MsgClaimTwitterAccount) Type() string {
	return MsgTypeClaimTwitterAccount
}

func (m MsgClaimTwitterAccount) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender addr")
	}
	if m.PreOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid prev owner addr")
	}
	if len(m.Handle) == 0 {
		return sdkerrors.Wrap(ErrInvalidInputField, "invalid handler")
	}
	return nil
}

func (m MsgClaimTwitterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgClaimTwitterAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

type MsgUpdateAccessList struct {
	Sender  sdk.AccAddress `json:"sender"`
	Address sdk.AccAddress `json:"address"`
}

func NewMsgUpdateAccessList(sender, address sdk.AccAddress) MsgUpdateAccessList {
	return MsgUpdateAccessList{
		Sender:  sender,
		Address: address,
	}
}

var _ sdk.Msg = MsgUpdateAccessList{}

func (m MsgUpdateAccessList) Route() string {
	return RouterKey
}

func (m MsgUpdateAccessList) Type() string {
	return "msg_update_access_list"
}

func (m MsgUpdateAccessList) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender addr")
	}
	
	if m.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid  addr ")
	}
	return nil
}

func (m MsgUpdateAccessList) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpdateAccessList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

type MsgUpdateHandlersInfo struct {
	Sender   sdk.AccAddress `json:"sender"`
	Handlers []string       `json:"handlers"`
}

func NewMsgUpdateHandlerInfo(sender sdk.AccAddress, handler []string) MsgUpdateHandlersInfo {
	return MsgUpdateHandlersInfo{
		Sender:   sender,
		Handlers: handler,
	}
}

var _ sdk.Msg = MsgUpdateHandlersInfo{}

func (m MsgUpdateHandlersInfo) Route() string {
	return RouterKey
}

func (m MsgUpdateHandlersInfo) Type() string {
	return "msg_updater_handlers_info"
}

func (m MsgUpdateHandlersInfo) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender addr")
	}
	if len(m.Handlers) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "length is zero")
	}
	return nil
}

func (m MsgUpdateHandlersInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgUpdateHandlersInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}
