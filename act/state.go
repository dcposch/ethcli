package act

import (
	"crypto/ecdsa"

	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Browser state
type State struct {
	Tab   TabState
	Chain ChainState
}

// Ethereum chain connection state
type ChainState struct {
	// logged-in account key
	PrivateKey *ecdsa.PrivateKey
	// logged-in account, eg vitalik.eth
	Account eth.NamedAddr
	// connection status, chain ID, etc
	Conn eth.ConnStatus
}

// Tab state
type TabState struct {
	// User entry in URL bar
	EnteredAddr string
	// Resolved contract addresss
	ContractAddr *common.Address
	// Error loading the app
	ErrorText string
	// Error within the app
	AppErrorText string
	// The displayed app, as returned by the contract render()
	Vdom []eth.VElem
	// ABI-encoded user inputs. Inputs[k] == nil if user hasn't entered anything for key k.
	Inputs [][]byte
	// Shows confirmation modal.
	ProposedTx *ethereum.CallMsg
	// Sent transaction, waiting for block confirmation.
	PendingTx *types.Transaction
}
