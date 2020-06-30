package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"time"
	
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
		Use:   types.ModuleName,
		Short: "IBC nfts  transfer transaction subcommands",
	}
	
	ics20XNFTTransferTxCmd.AddCommand(flags.PostCommands(
		GetXNFTTxCmd(cdc),
		GetSlotBookingTxCmd(cdc),
		GetSetParamsTxCmd(cdc),
		GetMsgPayLicensingFee(cdc),
		GetMsgDistributeFunds(cdc),
	)...)
	
	return ics20XNFTTransferTxCmd
}

// GetXNFTTxCmd returns the command to create a NewMsgTransfer transaction
func GetXNFTTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfts-transfer [src-port] [src-channel] [dest-height] [recipient]",
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

func GetSlotBookingTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book [src-port] [src-channel] [dest-height] [programTime] [liveStream-id]",
		Short: "Book Slot  through IBC",
		Args:  cobra.ExactArgs(5),
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
			
			primaryNFTID := viper.GetString(FlagPrimaryNFTID)
			adNFTID := viper.GetString(FlagAdNFTID)
			if primaryNFTID != "" && adNFTID != "" {
				return fmt.Errorf("cannot book slot with both ad and pNFT. Try with any one of it")
			}
			
			amountStr := viper.GetString(FlagAmount)
			
			if adNFTID != "" && amountStr == "" {
				return fmt.Errorf("must provide amount")
			}
			
			var coin sdk.Coin
			if amountStr != "" {
				coin, err = sdk.ParseCoin(amountStr)
				if err != nil {
					return err
				}
			}
			
			timeStamp, err := time.Parse(time.RFC3339, args[3])
			if err != nil {
				return err
				
			}
			
			if timeStamp.Before(time.Now().UTC()) {
				return fmt.Errorf("start date can't be the past date")
			}
			
			msg := types.NewMsgSlotBooking(primaryNFTID, adNFTID, args[4], sender, coin, timeStamp, srcPort, srcChannel, uint64(destHeight))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().AddFlagSet(PrimaryNFTID)
	cmd.Flags().AddFlagSet(AdNFTID)
	cmd.Flags().AddFlagSet(Amount)
	return cmd
}

func GetSetParamsTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params [coco-channel] [ff-channel] [dest-height] ",
		Short: "Set channel in param space",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			
			sender := cliCtx.GetFromAddress()
			nftChannel := args[0]
			
			destHeight, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			
			msg := types.NewMsgSetParams(sender, nftChannel, args[1], uint64(destHeight))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgPayLicensingFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-licensing-fee [src-port] [src-channel] [dest-height] [amount] [recipient] [primary-nft-id] ",
		Short: "This transaction is round trip tx, it will pay the licensing fee from coco account and get the secondary nfts from ff chain",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			coins, err := sdk.ParseCoin(args[3])
			if err != nil {
				return err
			}
			destHeight, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			
			msg := types.NewMsgPayLicensingFee(args[0], args[1], args[5], uint64(destHeight), coins, cliCtx.GetFromAddress(), args[4])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgDistributeFunds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distribute [src-channel] [dest-height]  ",
		Short: "distribute funds to completed programmes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			destHeight, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			
			msg := types.NewMsgDistributeFunds(args[0], uint64(destHeight), cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
