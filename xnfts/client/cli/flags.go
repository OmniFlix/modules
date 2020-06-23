package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPrimaryNFTID   = "pnft-id"
	FlagSecondaryNFTID = "snft-id"
	FlagAdNFTID        = "adnft_id"
	FlagAssetID        = "asset-id"
	FlagRecipient      = "recipient"
	
	FlagLicensingFee  = "licensing-fee"
	FlagRevenueShare  = "revenue-share"
	FlagTwitterHandle = "handle"
	FlagAmount        = "amount"
)

var (
	PrimaryNFTID   = flag.NewFlagSet("", flag.ContinueOnError)
	SecondaryNFTID = flag.NewFlagSet("", flag.ContinueOnError)
	AssetID        = flag.NewFlagSet("", flag.ContinueOnError)
	LicensingFee   = flag.NewFlagSet("", flag.ContinueOnError)
	RevenueShare   = flag.NewFlagSet("", flag.ContinueOnError)
	TwitterHandle  = flag.NewFlagSet("", flag.ContinueOnError)
	Recipient      = flag.NewFlagSet("", flag.ContinueOnError)
	AdNFTID        = flag.NewFlagSet("", flag.ContinueOnError)
	Amount         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	
	PrimaryNFTID.String(FlagPrimaryNFTID, "", "primary nft id")
	SecondaryNFTID.String(FlagSecondaryNFTID, "", "secondary nft id")
	Recipient.String(FlagRecipient, "", "receiver address")
	AdNFTID.String(FlagAdNFTID, "", "ad nft id")
	Amount.String(FlagAmount, "", "amount to be paid while adding ad nft")
}
