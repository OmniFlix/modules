package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNFTNotFound    = sdkerrors.Register(ModuleName, 11, "nfts not found")
	ErrInvalidLicense = sdkerrors.Register(ModuleName, 12, "invalid license")
	ErrParamsNotFound = sdkerrors.Register(ModuleName, 13, "params not found")
)
