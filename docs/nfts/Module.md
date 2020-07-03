
### Module

Our `scaffold` tool has done most of the work for us in generating our `module.go` file inside `./nfts/.` One way that our module is different than the simplest form of a module, is that it uses it's own `Keeper`. The only real changes needed are under the `AppModule` and `NewAppModule`. The file should look as follows afterward:

```go=
package nfts

import (
    "encoding/json"
    "fmt"
    
    "github.com/gorilla/mux"
    "github.com/spf13/cobra"
    abci "github.com/tendermint/tendermint/abci/types"
    
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"
    
    "github.com/FreeFlixMedia/modules/nfts/client/cli"
    "github.com/FreeFlixMedia/modules/nfts/internal/types"
)

var (
    _ module.AppModule      = AppModule{}
    _ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
    types.RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
    return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
    var data types.GenesisState
    if err := cdc.UnmarshalJSON(bz, &data); err != nil {
        return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
    }
    
    return data.ValidateGenesis()
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetQueryCmd(cdc)
}

type AppModule struct {
    AppModuleBasic
    nftKeeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
    return AppModule{
        nftKeeper: keeper,
    }
}

func (AppModule) Name() string {
    return ModuleName
}

func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() sdk.Handler { return NewHandler(am.nftKeeper) }

func (AppModule) QuerierRoute() string {
    return types.QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
    return NewQuerier(am.nftKeeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
    var genesisState GenesisState
    cdc.MustUnmarshalJSON(data, &genesisState)
    InitGenesis(ctx, am.nftKeeper, genesisState)
    return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
    gs := ExportGenesis(ctx, am.nftKeeper)
    return cdc.MustMarshalJSON(gs)
}

func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
    return []abci.ValidatorUpdate{}
}

```

Congratulations you have completed the `nfts` module!

This module is now able to be incorporated into any Cosmos SDK application.

Since we don't want to just build a module but want to build an application that also uses that module, let's go through the process of configuring an app.
