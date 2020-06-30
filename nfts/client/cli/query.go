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
		Short:                      "Querying commands for the nfts module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	
	cmd.AddCommand(
		GetCmdQueryTweetNFT(cdc),
		GetCmdQueryTweetsByAccount(cdc),
		GetCmdQueryLiveStream(cdc),
		GetCmdQueryDNFT(cdc),
		GetCmdQueryLiveStreams(cdc),
		GetCmdQueryAdNFT(cdc),
		GetCmdQueryAdsByAccount(cdc),
		GetCmdQueryTwitterAccountInfo(cdc),
		GetCmdQueryTwitterAccounts(cdc),
		GetCmdQueryAuthorisedHandlers(cdc),
		GetCmdQueryAuthorisedAddresses(cdc),
	)
	
	return cmd
}

func GetCmdQueryTweetNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfts [id]",
		Short: "Get NFT using nfts id ",
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

func GetCmdQueryLiveStream(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liveStream [id]",
		Short: "Get liveStream using liveStream id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryLiveStream, args[0]), nil)
			if err != nil {
				return err
			}
			
			var liveStream types.BaseLiveStream
			cdc.MustUnmarshalJSON(bz, &liveStream)
			return cliCtx.PrintOutput(liveStream)
		},
	}
	
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryDNFT(cdc *codec.Codec) *cobra.Command { // TODO : Query DNFT with DNFT ID
	cmd := &cobra.Command{
		Use:   "dnft [slotTime]",
		Short: "Query dnft by slot time",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryDNFT, args[0]), nil)
			if err != nil {
				return err
			}
			
			var dnft types.BaseDNFT
			cdc.MustUnmarshalJSON(bz, &dnft)
			return cliCtx.PrintOutput(dnft)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryLiveStreams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liveStreams",
		Short: "Query  live streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiveStreams), nil)
			if err != nil {
				return err
			}
			
			var streams []types.BaseLiveStream
			cdc.MustUnmarshalJSON(bz, &streams)
			return cliCtx.PrintOutput(streams)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryAdNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adnft [id]",
		Short: "Query ad-nfts by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAdNFTBYID, args[0]), nil)
			if err != nil {
				return err
			}
			
			var adnft types.BaseAdNFT
			cdc.MustUnmarshalJSON(bz, &adnft)
			return cliCtx.PrintOutput(adnft)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryAdsByAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adnfts [address]",
		Short: "Get Ad NFTs associated to account address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAdNFTByAddress, args[0]), nil)
			if err != nil {
				return err
			}
			
			var adNFTs []types.BaseAdNFT
			cdc.MustUnmarshalJSON(res, &adNFTs)
			return cliCtx.PrintOutput(adNFTs)
		},
	}
	return flags.GetCommands(cmd)[0]
	
}

func GetCmdQueryTwitterAccountInfo(cdc *codec.Codec) *cobra.Command { // TODO : Query DNFT with DNFT ID
	cmd := &cobra.Command{
		Use:   "account [twitter-handle]",
		Short: "Query Twitter account associated twitter account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTwitterAccount, args[0]), nil)
			if err != nil {
				return err
			}
			
			var info types.TwitterAccountInfo
			cdc.MustUnmarshalJSON(bz, &info)
			return cliCtx.PrintOutput(info)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryTwitterAccounts(cdc *codec.Codec) *cobra.Command { // TODO : Query DNFT with DNFT ID
	cmd := &cobra.Command{
		Use:   "accounts ",
		Short: "Query Twitter accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTwitterAccounts, nil), nil)
			if err != nil {
				return err
			}
			
			var info []types.TwitterAccountInfo
			cdc.MustUnmarshalJSON(bz, &info)
			return cliCtx.PrintOutput(info)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryAuthorisedHandlers(cdc *codec.Codec) *cobra.Command { // TODO : Query DNFT with DNFT ID
	cmd := &cobra.Command{
		Use:   "acl-handlers ",
		Short: "Query authorised handlers to mint nfts ",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAuthorisedHandler, nil), nil)
			if err != nil {
				return err
			}
			
			var info types.AllowedHandles
			cdc.MustUnmarshalJSON(bz, &info)
			return cliCtx.PrintOutput(info)
		},
	}
	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryAuthorisedAddresses(cdc *codec.Codec) *cobra.Command { // TODO : Query DNFT with DNFT ID
	cmd := &cobra.Command{
		Use:   "acl-addresses ",
		Short: "Query authorised addresses to upate handler ",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAuthorisedAddresses, nil), nil)
			if err != nil {
				return err
			}
			
			var info types.AclInfo
			cdc.MustUnmarshalJSON(bz, &info)
			return cliCtx.PrintOutput(info)
		},
	}
	return flags.GetCommands(cmd)[0]
}
