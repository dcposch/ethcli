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
	// logged-in account, eg vitalik.eth
	Account eth.NamedAddr
	// connection status, chain ID, etc
	Conn eth.ConnStatus
}

// Tab state
type TabState struct {
	EnteredAddr  string
	ContractAddr *common.Address
	ErrorText    string
	Vdom         []eth.VdomElem
}
