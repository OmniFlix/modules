package xnfts

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/FreeFlixMedia/modules/nfts"
)

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	
	liveStreams := keeper.GetAllLiveStreams(ctx)
	for _, stream := range liveStreams {
		payout := stream.Payout
		endTime := ctx.BlockTime().UTC()
		endSlotTime := endTime.Format("2006-01-0215:04")
		
		// get dnfts of payout time + more 5 min
		duration := payout + (5 * time.Minute)
		startTime := endTime.Add(-duration)
		startSlotTime := startTime.Format("2006-01-0215-04")
		fmt.Println("TIME SLOTS ========================", startSlotTime, endSlotTime)
		baseDnfts := keeper.GetDNFTsBetweenInterval(ctx, startSlotTime, endSlotTime)
		
		for _, baseDNFT := range baseDnfts {
			var amountPaidToPrimaryOwner sdk.Coin
			var primaryNFTOwner string
			
			if baseDNFT.Status == "UNPAID" {
				
				fmt.Println("STREAM_REVENUE_SHARE", stream.RevenueShare)
				amount := stream.RevenueShare.MulInt(baseDNFT.LockedAmount.Amount)
				
				if stream.RevenueShare.RoundInt64() > sdk.OneDec().RoundInt64() {
					continue
				}
				
				if len(baseDNFT.AdNFTID) < 1 {
					continue
				}
				
				amountPaidtoLiveStream := sdk.NewInt64Coin(baseDNFT.LockedAmount.Denom, amount.RoundInt64())
				baseDNFT.LockedAmount = baseDNFT.LockedAmount.Sub(amountPaidtoLiveStream)
				_, err := keeper.AddCoins(ctx, stream.OwnerAddress, sdk.Coins{amountPaidtoLiveStream})
				if err != nil {
					continue
				}
				
				if baseDNFT.Type == nfts.TypeLocal {
					secondaryNFT, _ := keeper.GetTweetNFTByID(ctx, baseDNFT.NFTID)
					
					fmt.Println("SECONDARY_NFT_REVENUE_SHARE", secondaryNFT.RevenueShare)
					if secondaryNFT.RevenueShare.RoundInt64() > sdk.OneDec().RoundInt64() {
						continue
					}
					amount = secondaryNFT.RevenueShare.MulInt(baseDNFT.LockedAmount.Amount)
					amountPaidToPrimaryOwner = sdk.NewInt64Coin(baseDNFT.LockedAmount.Denom, amount.RoundInt64())
					baseDNFT.LockedAmount = baseDNFT.LockedAmount.Sub(amountPaidToPrimaryOwner)
					primaryNFTOwner = secondaryNFT.PrimaryOwner
					
					amountPaidToSecondaryOwner := baseDNFT.LockedAmount
					_, err := keeper.AddCoins(ctx, GetHexAddressFromBech32String(secondaryNFT.SecondaryOwner), sdk.Coins{amountPaidToSecondaryOwner})
					if err != nil {
						continue
					}
					baseDNFT.LockedAmount = baseDNFT.LockedAmount.Sub(amountPaidToSecondaryOwner)
				}
				if baseDNFT.Type == nfts.TypeIBC {
					amountPaidToPrimaryOwner = baseDNFT.LockedAmount
					primaryNFTOwner = baseDNFT.PrimaryNFTAddress
					baseDNFT.LockedAmount = baseDNFT.LockedAmount.Sub(amountPaidToPrimaryOwner)
				}
				
				// IBC Transfer
				packet := PacketTokenDistribution{
					Recipient:    primaryNFTOwner,
					Handler:      baseDNFT.TwitterHandleName,
					AmountLocked: amountPaidToPrimaryOwner,
				}
				
				params, _ := keeper.GetParams(ctx) // TODO automate params to destination client
				
				if err := keeper.XTransfer(ctx, XNFTPortID, params.NFTChannel, params.DestHeight, packet.GetBytes()); err != nil {
					baseDNFT.LockedAmount = baseDNFT.LockedAmount.Add(amountPaidToPrimaryOwner)
				}
				
				baseDNFT.Status = "PAID"
				keeper.SetDNFT(ctx, baseDNFT)
				ctx.EventManager().EmitEvents(
					sdk.Events{sdk.NewEvent(
						nfts.EventTypeDistribution,
						sdk.NewAttribute(nfts.AttributeDNFTID, baseDNFT.DNFTID),
						sdk.NewAttribute(nfts.AttributeLiveStreamID, baseDNFT.LiveStreamID),
					),
						sdk.NewEvent(
							EventTypeTokenDistribution,
							sdk.NewAttribute(nfts.AttributeDNFTID, baseDNFT.DNFTID),
							sdk.NewAttribute(AttributeKeyReceiver, packet.Recipient),
						),
					},
				)
			}
		}
	}
	
}
