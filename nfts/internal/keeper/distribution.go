package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error) {
	return k.bankKeeper.AddCoins(ctx, addr, amount)
}
