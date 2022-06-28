package act

import (
	"log"
	"math/big"
	"strings"

	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum/common"
)

// Represents a single user action.
// Dispatched from the UI, handled by act.Dispatch()
type Action interface {
	Run()
}

// Navigating, either via link or via the URL bar.
type ActSetUrl struct {
	Url string
}

func (a *ActSetUrl) Run() {
	url := a.Url
	tab := &state.Tab

	tab.EnteredAddr = url
	tab.ErrorText = ""
	tab.ContractAddr = nil

	// Navigation is always to a UI contract.
	// User either enters an address directly, or an ENS name.
	if strings.HasSuffix(url, ".eth") {
		render()
		result, err := client.Resolve(url)
		if err != nil {
			tab.ErrorText = err.Error()
		} else {
			tab.ContractAddr = &result
		}
	} else if strings.HasPrefix(url, "0x") {
		addr := common.HexToAddress(url)
		tab.ContractAddr = &addr
	} else {
		tab.EnteredAddr = ""
	}
	tab.Inputs = make(map[string][]byte)

	reloadTab()
}

// Update a form input.
type ActSetInput struct {
	Key big.Int
	Val []byte
}

func (a *ActSetInput) Run() {
	state.Tab.Inputs[a.Key.String()] = a.Val
}

// Submit a form.
type ActSubmit struct {
	ButtonKey big.Int
}

func (a *ActSubmit) Run() {
	// TODO: construct act() transaction
	log.Printf("act Submit %s %#v", a.ButtonKey.String(), state.Tab.Inputs)
}

// Reload context information about the blockchain.
func reloadChainState() {
	state.Chain.Conn = client.ConnStatus()
	state.Chain.Account = eth.NamedAddr{}
}

func reloadTab() {
	if state.Tab.ContractAddr == nil {
		return
	}

	appState := []byte{}
	vdom, err := client.FrontendRender(
		state.Chain.Account.Addr, *state.Tab.ContractAddr, appState)
	if err == nil {
		// TODO: vdom diffing
		state.Tab.Vdom = vdom
		state.Tab.ErrorText = ""
		for _, v := range vdom {
			key := v.DataElem.(eth.KeyElem).GetKey()
			state.Tab.Inputs[key.String()] = []byte{}
		}
	} else {
		state.Tab.Vdom = nil
		state.Tab.ErrorText = err.Error()
	}

	render()
}

func render() {
	renderer(&state)
}
