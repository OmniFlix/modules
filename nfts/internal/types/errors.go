package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNFTNotFound         = sdkerrors.Register(ModuleName, 11, "nfts not found")
	ErrInvalidLicense      = sdkerrors.Register(ModuleName, 12, "invalid license")
	ErrInvalidInputField   = sdkerrors.Register(ModuleName, 13, "invalid field")
	ErrAccountNotFound     = sdkerrors.Register(ModuleName, 14, "account doesn't exist")
	ErrInvalidClaimStatus  = sdkerrors.Register(ModuleName, 15, "claim status invalid")
	ErrAssetIDAlreadyExist = sdkerrors.Register(ModuleName, 16, "asset id already exist")
	ErrInvalidSlotBooking  = sdkerrors.Register(ModuleName, 17, "slot booking")
	ErrAddressAlreadyExist = sdkerrors.Register(ModuleName, 18, "address already exist")
	ErrHandlerAlreadyExist = sdkerrors.Register(ModuleName, 19, "handler already exist")
)
