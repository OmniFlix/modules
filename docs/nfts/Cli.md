

### CLI

A Command Line Interface (CLI) will help us interact with our app once it is running on a machine somewhere. Each Module has its own namespace within the CLI that gives it the ability to create and sign Messages destined to be handled by that module. It also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI will let you do things like generate keys for a new account or check the status of an interaction you already had with the application.


The CLI for our module is broken into two files called `tx.go` and `query.go` which are located in `./nfts/client/cli/.` One file is for making transactions that contain messages which will ultimately update our state. The other is for making queries which will give us the ability to read information from our state. Both files utilize the [Cobra](https://github.com/spf13/cobra) library.

**Transactions**

The `tx.go` file contains `GetTxCmd` which is a standard method within the Cosmos SDK. It is referenced later in the `module.go` file which describes exactly which attributes a module has. This makes it easier to incorporate different modules for different reasons at the level of the actual application. After all, we are focusing on a module at this point, but later we will create an application that utilizes this module as well as other modules that are already available within the Cosmos SDK.

Inside `GetTxCmd` we create a new module-specific command and call is `nfts`. Within this command we add a sub-command for each Message type we've defined:

* `GetMsgMintTweetNFT`


Each function takes parameters from the **Cobra** CLI tool to create a new msg, sign it and submit it to the application to be processed. These functions should go into the `tx.go` file and look as follows:

```go=
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

```



**Query**

The `query.go` file contains similar **Cobra** commands that reserve a new namespace for referencing our `nfts` module. Instead of creating and submitting messages, however, the `query.go` the file creates queries and returns the results in human-readable form. The queries it handles are the same we defined in our `querier.go` file earlier:

* `GetCmdQueryTweetNFT`
* `GetCmdQueryTweetsByAccount`

After defining these commands, your `query.go` file should look like:

```go=
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

```

While these are all the major moving pieces of a module (`Message`, `Handler`, `Keeper`, `Querier`, and `Client`) there are some organizational tasks that we have yet to complete. The next step will be making sure that our module is completely configured to make it usable within any application.
