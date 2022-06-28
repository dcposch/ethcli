package act

import (
	"log"
	"strings"

	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	tab.AppErrorText = ""

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

	reloadTab()
}

// Update a form input.
type ActSetInput struct {
	Key uint8
	Val []byte
}

func (a *ActSetInput) Run() {
	state.Tab.Inputs[a.Key] = a.Val
}

// Submit a form.
type ActSubmit struct {
	ButtonKey uint8
}

func (a *ActSubmit) Run() {
	inputs := state.Tab.Inputs
	log.Printf("act Submit %d", a.ButtonKey)

	appState := []byte{}
	contractAddr := *state.Tab.ContractAddr
	action := eth.ButtonAction{ButtonKey: a.ButtonKey, Inputs: inputs}
	newAppState, err := client.FrontendSubmit(state.Chain.Account.Addr, contractAddr, appState, action)
	log.Printf("act Submit %d result %v err %v", a.ButtonKey, hexutil.Encode(newAppState), err)

	if err == nil {
		state.Tab.AppErrorText = ""
	} else {
		state.Tab.AppErrorText = err.Error()
	}

	render()
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
		state.Tab.Vdom = vdom
		state.Tab.ErrorText = ""

		maxId := uint8(0)
		for _, v := range state.Tab.Vdom {
			if v.DataElem.GetKey() > maxId {
				maxId = v.DataElem.GetKey()
			}
		}
		state.Tab.Inputs = make([][]byte, maxId+1)
	} else {
		state.Tab.Vdom = nil
		state.Tab.ErrorText = err.Error()
	}

	render()
}

func render() {
	renderer(&state)
}
