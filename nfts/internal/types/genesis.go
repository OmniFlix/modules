package types

type GenesisState struct {
	TweetNFTs          []BaseTweetNFT       `json:"tweet_nfts"`
	AdNFTs             []BaseAdNFT          `json:"ad_nft"`
	DNFTs              []BaseDNFT           `json:"dnf_ts"`
	LiveStreams        []BaseLiveStream     `json:"live_streams"`
	TwitterAccountInfo []TwitterAccountInfo `json:"twitter_account_info"`
	ACLAddressList     AclInfo              `json:"acl_address_list"`
	ACLHandlersInfo    AllowedHandles       `json:"acl_handlers_info"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs GenesisState) ValidateGenesis() error {
	return nil // TODO Validate
}
