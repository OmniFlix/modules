### Events

At the end of each handler is an EventManager which will create logs within the transaction that reveals information about what occurred during the handling of this message. This is useful for client-side software that wants to know exactly what happened as a result of this state transition. These Events use a series of pre-defined types that can be found in `./xnfts/internal/types/events.go` and look as follows:

```go=
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


```

Now that we have all the necessary pieces for updating state (`Message`, `Handler`, `Keeper`) we might want to consider ways in which we can query state. This is typically done via a REST endpoint and/or a CLI. Both of those clients interact with a part of the app which queries state, called the `Querier`.