package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNFTNotFound         = sdkerrors.Register(ModuleName, 11, "nft not found")
	ErrAssetIDAlreadyExist = sdkerrors.Register(ModuleName, 12, "asset id already exist")
	
	ErrInvalidLicense = sdkerrors.Register(ModuleName, 13, "invalid license")
	ErrParamsNotFound = sdkerrors.Register(ModuleName, 14, "params not found")
)
