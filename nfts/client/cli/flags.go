package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagLicence       = "licence"
	FlagLicenceFee    = "licence-fee"
	FlagRevenueShare  = "revenue-share"
	FlagTwitterHandle = "handle"
	FlagAssetID       = "asset-id"
)

var (
	fsLicence      = flag.NewFlagSet("", flag.ContinueOnError)
	fsLicenceFee   = flag.NewFlagSet("", flag.ContinueOnError)
	fsRevenueShare = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsLicence.Bool(FlagLicence, false, "To give licence for user")
	fsLicenceFee.String(FlagLicenceFee, "", "Amount to pay to get licence for nfts")
	fsRevenueShare.String(FlagRevenueShare, "", "Percentage of amount to get after giving licence")
}
