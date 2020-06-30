package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
	
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	NFTTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "nfts  transfer transaction subcommands",
	}
	
	NFTTxCmd.AddCommand(flags.PostCommands(
		GetMsgCreateInitialTweetNFT(cdc),
		GetMsgMintTweetNFT(cdc),
		GetMsgCreateDNFT(cdc),
		GetMsgCreateAdNFT(cdc),
		GetMsgCreateLiveStream(cdc),
		GetMsgAccessAddressList(cdc),
		GetMsgUpdateHandler(cdc),
	)...)
	
	return NFTTxCmd
}

func GetMsgCreateInitialTweetNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-nfts [asset-id] [twitter-handle]",
		Short: "mint initial nfts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			msg := types.NewMsgCreateInitialTweetNFT(cliCtx.GetFromAddress(), args[0], args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	
	return cmd
}

func GetMsgMintTweetNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nfts ",
		Short: "mint initial nfts",
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

func GetMsgCreateLiveStream(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-liveStream [costPerAdPerSlot] [revenue-share] [payoutInMin] [slotsPerMin]",
		Short: "live_stream creation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inbuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inbuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inbuf).WithCodec(cdc)
			
			cost, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}
			
			revenueShare, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}
			payout, err := time.ParseDuration(args[2])
			if err != nil {
				return err
			}
			
			slots, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}
			
			msg := types.NewMsgLiveStream(cliCtx.GetFromAddress(), revenueShare, cost, payout, uint64(slots))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgCreateDNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-dnft [sNFTID] [programTime] [live_streamID]",
		Short: "Booking slot in Livestream",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			programmeTime, err := time.Parse(time.RFC3339, args[1]) // RFC3339 date ex: 2020-03-03T06:26:19.862851614Z
			if err != nil {
				return err
			}
			
			if programmeTime.Before(time.Now().UTC()) {
				return fmt.Errorf("programme's start date can't be in the past")
			}
			
			fromAddress := cliCtx.GetFromAddress()
			
			msg := types.NewMsgBookSlot(fromAddress, args[0], args[2], programmeTime)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgCreateAdNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-adnft [assetid]",
		Short: "Booking slot in Livestream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			msg := types.NewMsgCreateAdNFT(args[0], cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgAccessAddressList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl [address]",
		Short: "Update the Acl Access",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			
			msg := types.NewMsgUpdateAccessList(cliCtx.GetFromAddress(), address)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetMsgUpdateHandler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-handle [handler1,handle2]",
		Short: "Add handler to access list of handlers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			
			handlersStr := strings.TrimSpace(args[0])
			if len(handlersStr) == 0 {
				return fmt.Errorf("handlers are empty")
			}
			
			var handler []string
			data := strings.Split(handlersStr, ",")
			for _, handle := range data {
				match := types.ReHandle.FindStringSubmatch(handle)
				if match == nil {
					return fmt.Errorf("invalid handler expression")
				}
				handler = append(handler, handle)
			}
			msg := types.NewMsgUpdateHandlerInfo(cliCtx.GetFromAddress(), handler)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
