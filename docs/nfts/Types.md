
### Types
Our queries are rather simple since we've already outfitted our `Keeper` with all the necessary functions to access the state. You can see the iterator being used here as well.

Now that we have all of the basic actions of our module created, we want to make them accessible. We can do this with a CLI client and a REST client. For this tutorial, we will just be creating a CLI client. If you are interested in what goes into making a REST client.

Let's take a look at what goes into making a CLI.
#### BaseTweetNFT

You may notice the reference to `types.BaseTweetNFT` throughout the `Keeper`. These is new struct defined in `./nfts/types/nfts.go` that contains all necessary information about different tweet nfts. You can create this file now and add the following:

```go=
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

```

You might also notice that each type has the `String` method. This allows us to render the struct as a string for rendering.