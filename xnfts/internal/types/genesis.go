package types

type GenesisState struct {
	PortID string `json:"port_id"`
}

func DefaultGenesis() GenesisState {
	return GenesisState{PortID: PortID}
}
