package cli

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	
	"github.com/FreeFlixMedia/modules/nfts/internal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the nft module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	
	cmd.AddCommand(
		GetCmdQueryTweetNFT(cdc),
		GetCmdQueryTweetsByAccount(cdc),
	)
	
	return cmd
}

func GetCmdQueryTweetNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft [id]",
		Short: "Get NFT using nft id ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTweetNFT, args[0]), nil)
			if err != nil {
				return err
			}
			
			var tweetNFT types.BaseTweetNFT
			cdc.MustUnmarshalJSON(res, &tweetNFT)
			return cliCtx.PrintOutput(tweetNFT)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryTweetsByAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfts [address]",
		Short: "Get  NFTs associated to account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTweetNFTsByAddress, args[0]), nil)
			if err != nil {
				return err
			}
			
			var tweetNFTs []types.BaseTweetNFT
			cdc.MustUnmarshalJSON(res, &tweetNFTs)
			return cliCtx.PrintOutput(tweetNFTs)
		},
	}
	return flags.GetCommands(cmd)[0]
	
}
