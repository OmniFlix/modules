package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BaseTweetNFT struct {
	PrimaryNFTID string `json:"primary_nft_id"`
	PrimaryOwner string `json:"primary_owner"`
	
	SecondaryNFTID string `json:"secondary_nft_id"`
	SecondaryOwner string `json:"secondary_owner"`
	
	License bool   `json:"license"`
	AssetID string `json:"asset_id"`
	
	LicensingFee sdk.Coin `json:"licensing_fee"`
	RevenueShare sdk.Dec  `json:"revenue_share"`
	
	TwitterHandle string `json:"twitter_handle"`
}

func (nft BaseTweetNFT) String() string {
	return fmt.Sprintf(`
PrimaryNFTID: %s,
PrimaryOwner: %s,

SecondaryNFTID: %s,
SecondaryOwner: %s,

License: %t,
AssetID: %s,

LicensingFee: %s,
RevenueShare: %s,

TwitterHandle: %s,
`, nft.PrimaryNFTID, nft.PrimaryOwner, nft.SecondaryNFTID, nft.SecondaryOwner,
		nft.License, nft.AssetID, nft.LicensingFee.String(), nft.RevenueShare.String(), nft.TwitterHandle)
}

type BaseAdNFT struct {
	AdNFTID string `json:"ad_nftid"`
	Owner   string `json:"owner"`
	AssetID string `json:"asset_id"`
}

func (nft BaseAdNFT) String() string {
	return fmt.Sprintf(`
AdNFTID: %s,
Owner: %s,
AssetID:
`, nft.AdNFTID, nft.Owner, nft.AssetID)
}

type TwitterAccountInfo struct {
	Owner        sdk.AccAddress `json:"owner"`
	Handle       string         `json:"handle"`
	ClaimStatus  bool           `json:"claim_status"`
	LockedAmount sdk.Coins      `json:"locked_amount"`
}

func (info TwitterAccountInfo) String() string {
	return fmt.Sprintf(`
Owner: %s,
Handle: %s,
ClaimStatus: %t,
LockedAmount: %s,
`, info.Owner.String(), info.Handle, info.ClaimStatus, info.LockedAmount.String())
}

type AclInfo struct {
	AccessList []sdk.AccAddress `json:"access_list"`
}

func (acl AclInfo) String() string {
	if len(acl.AccessList) == 0 {
		return ""
	}
	
	out := ""
	for _, addr := range acl.AccessList {
		out += fmt.Sprintf("%v,", addr.String())
	}
	return out[:len(out)-1]
}

type AllowedHandles struct {
	Handles []string `json:"handles"`
}

func (handles AllowedHandles) String() string {
	if len(handles.Handles) == 0 {
		return ""
	}
	
	out := ""
	for _, handle := range handles.Handles {
		out += fmt.Sprintf("%v,", handle)
	}
	return out[:len(out)-1]
}
