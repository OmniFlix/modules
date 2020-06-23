package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
	return "msg_mint_tweet_nft"
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
