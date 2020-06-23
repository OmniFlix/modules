package xnfts

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/FreeFlixMedia/modules/xnfts/internal/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, state types.GenesisState) {
	if !keeper.IsBounded(ctx, state.PortID) {
		err := keeper.BindPort(ctx, state.PortID)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	portID := keeper.GetPort(ctx)
	
	return types.GenesisState{
		PortID: portID,
	}
}
