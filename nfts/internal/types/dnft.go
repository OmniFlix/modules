package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeLocal = "local"
	TypeIBC   = "ibc"
)

type BaseDNFT struct {
	DNFTID            string    `json:"dnft_id"`
	ProgramTime       time.Time `json:"program_time"` // RFC3339 date ex: 2020-03-03T06:26:19.862851614Z
	AdNFTID           string    `json:"ad_nft_id"`
	AdNFTAssetID      string    `json:"ad_nft_asset_id"`
	NFTID             string    `json:"nft_id"`
	TweetAssetID      string    `json:"tweet_asset_id"`
	PrimaryNFTAddress string    `json:"primary_nft_address"`
	TwitterHandleName string    `json:"twitter_handle_name"`
	Status            string    `json:"status"`
	LiveStreamID      string    `json:"live_stream_id"`
	Type              string    `json:"type"`
	LockedAmount      sdk.Coin  `json:"locked_amount"`
}

func GenerateDNFTID(programTime time.Time) string {
	return string(sdk.FormatTimeBytes(programTime.UTC()))
}

func GetBech32StringOfSlotAndProgrammeDetails(programTime time.Time) (string, string) {
	slotTime := programTime.Format("2006-01-0215:04")
	programmeTime := programTime.Format("2006-01-0215:04:05")
	return slotTime, programmeTime
	
}

func GetTimeSlotFromDNFTID(id string) (string, string) {
	timeC, _ := sdk.ParseTimeBytes([]byte(id))
	return GetBech32StringOfSlotAndProgrammeDetails(timeC)
}

func NewBaseNFT(programTime string) BaseDNFT {
	return BaseDNFT{
		DNFTID: programTime,
	}
}
