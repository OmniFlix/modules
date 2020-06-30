package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	StatusActive   = "active"
	StatusInActive = "inactive"
)

type BaseLiveStream struct {
	OwnerAddress     sdk.AccAddress `json:"owner_address"`
	LiveStreamID     string         `json:"live_stream_id"`
	CostPerAdPerSlot sdk.Coin       `json:"cost_per_ad_per_slot"`
	RevenueShare     sdk.Dec        `json:"revenue_share"`
	DNFTIDs          []string       `json:"dnft_ids"`
	Payout           time.Duration  `json:"payout"`
	SlotsPerMinute   uint64         `json:"slots_per_minute"`
	Status           string         `json:"status"`
}
