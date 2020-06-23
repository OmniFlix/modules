package types

type GenesisState struct {
	TweetNFTs []BaseTweetNFT `json:"tweet_nfts"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs GenesisState) ValidateGenesis() error {
	return nil // TODO Validate
}
