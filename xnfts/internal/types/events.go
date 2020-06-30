package types

import (
	"fmt"
	
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

const (
	EventTypeNFTPacketTransfer             = "nft_packet_transfer"
	EventTypeMintNFT                       = "mint_nft"
	EventTypeUpdateNFT                     = "update_nft"
	EventTypeSetParams                     = "set_params"
	EventTypeTokenDistribution             = "token_distribution"
	EventTypePayLicensingFeeAndNFTTransfer = "pay_licensing_fee_and_token_transfer"
	
	AttributeKeyReceiver    = "receiver"
	AtttibutePrimaryNFTID   = "primary_nft_id"
	AtttibuteSecondaryNFTID = "seconday_nft_id"
	AttributeKeyAckSuccess  = "success"
	AttributeKeyAckError    = "error"
	AttributeKeyNFTChannel  = "nft_channel"
	AttributeKeyDestHeight  = "dest_height"
)

var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)
