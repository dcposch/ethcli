package act

import (
	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum/common"
)

// Browser state
type State struct {
	Tab   TabState
	Chain ChainState
}

// Ethereum chain connection state
type ChainState struct {
	// 1 for mainnet, etc
	ChainID int64
	// logged-in account, eg vitalik.eth
	Account eth.NamedAddr
}

// Tab state
type TabState struct {
	EnteredAddr  string
	ContractAddr *common.Address
	ErrorText    string
}
