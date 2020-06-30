package types

type Params struct {
	NFTChannel string `json:"nft_channel"`
	FFChannel  string `json:"ff_channel"`
	DestHeight uint64 `json:"dest_height"`
}

func NewParams(cocoChannel, ffChannel string, destHeight uint64) Params {
	return Params{
		NFTChannel: cocoChannel,
		FFChannel:  ffChannel,
		DestHeight: destHeight,
	}
}
