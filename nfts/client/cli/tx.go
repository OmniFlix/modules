package cli

import (
	"bufio"
	"strconv"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	NFTTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "nfts  transfer transaction subcommands",
	}
	
	NFTTxCmd.AddCommand(flags.PostCommands(
		GetMsgMintTweetNFT(cdc),
	)...)
	
	return NFTTxCmd
}

func GetMsgMintTweetNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft ",
		Short: "mint tweet nft",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			var fee sdk.Coin
			var share sdk.Dec
			var err error
			
			license, err := strconv.ParseBool(viper.GetString(FlagLicence))
			if err != nil {
				return err
			}
			
			feeStr := viper.GetString(FlagLicenceFee)
			if feeStr != "" {
				fee, err = sdk.ParseCoin(feeStr)
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
			
			msg := types.NewMsgMintNFT(cliCtx.GetFromAddress(), viper.GetString(FlagAssetID), license, fee, share, viper.GetString(FlagTwitterHandle))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().String(FlagAssetID, "", "AssetID")
	cmd.Flags().String(FlagTwitterHandle, "", "Twitter handle")
	cmd.Flags().String(FlagLicenceFee, "0coco", "Twitter handle")
	cmd.Flags().String(FlagRevenueShare, "0", "Revenue share")
	cmd.Flags().String(FlagLicence, "false", "license")
	return cmd
}
