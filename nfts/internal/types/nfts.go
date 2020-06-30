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
