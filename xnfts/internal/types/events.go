package types

import (
	"fmt"
	
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

var (
	EventTypeNFTPacketTransfer             = "nft_packet_transfer"
	EventTypePayLicensingFeeAndNFTTransfer = "pay_licensing_fee_and_token_transfer"
	
	AttributeKeyReceiver   = "receiver"
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)
