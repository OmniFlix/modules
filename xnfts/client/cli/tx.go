package cli

import (
	"bufio"
	"strconv"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

// GetTxCmd returns the transaction commands for IBC fungible token transfer
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	ics20XNFTTransferTxCmd := &cobra.Command{
		Use:   "xnfts",
		Short: "IBC nft  transfer transaction subcommands",
	}
	
	ics20XNFTTransferTxCmd.AddCommand(flags.PostCommands(
		GetXNFTTxCmd(cdc),
	)...)
	
	return ics20XNFTTransferTxCmd
}

// GetXNFTTxCmd returns the command to create a NewMsgTransfer transaction
func GetXNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-transfer [src-port] [src-channel] [dest-height] [recipient]",
		Short: "Transfer non fungible token through IBC",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			
			sender := cliCtx.GetFromAddress()
			srcPort := args[0]
			srcChannel := args[1]
			destHeight, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			
			var fee sdk.Coin
			var share sdk.Dec
			var assteID, handle string
			var data types.NFTInput
			
			licenceFee := viper.GetString(FlagLicensingFee)
			if licenceFee != "" {
				fee, err = sdk.ParseCoin(licenceFee)
				if err != nil {
					return err
				}
			}
			
			shareStr := viper.GetString(FlagRevenueShare)
			if shareStr != "" {
				share, err = sdk.NewDecFromStr(shareStr)
				if err != nil {
					return err
				}
				
			}
			
			assteID = viper.GetString(FlagAssetID)
			handle = viper.GetString(FlagTwitterHandle)
			
			data = types.NFTInput{
				PrimaryNFTID:  viper.GetString(FlagPrimaryNFTID),
				Recipient:     args[3],
				AssetID:       assteID,
				LicensingFee:  fee,
				RevenueShare:  share,
				TwitterHandle: handle,
			}
			
			msg := types.NewMsgXNFTTransfer(srcPort, srcChannel, uint64(destHeight), sender, data)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().String(FlagPrimaryNFTID, "", "Primary NFT id")
	cmd.Flags().String(FlagRevenueShare, "0", "Revenue share")
	cmd.Flags().String(FlagLicensingFee, "0coco", "Licenese fee")
	cmd.Flags().String(FlagAssetID, "", "AssetID")
	cmd.Flags().String(FlagTwitterHandle, "", "Twitter Handle")
	return cmd
}
