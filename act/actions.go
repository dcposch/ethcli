package act

import (
	"context"
	"fmt"
	"log"
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
	appState := []byte{}
	contractAddr := *state.Tab.ContractAddr
	action := eth.ButtonAction{ButtonKey: a.ButtonKey, Inputs: inputs}
	callMsg, err := client.FrontendSubmit(state.Chain.Account.Addr, contractAddr, appState, action)
	log.Printf("act Submit %d err %v", a.ButtonKey, err)

	if err == nil {
		state.Tab.AppErrorText = ""
		state.Tab.ProposedTx = callMsg
	} else {
		state.Tab.AppErrorText = err.Error()
	}

	render()
}

type ActExecTx struct {
}

func (a *ActExecTx) Run() {
	tx, err := client.Execute(state.Tab.ProposedTx, state.Chain.PrivateKey)
	state.Tab.ProposedTx = nil
	if err == nil {
		state.Tab.PendingTx = tx
	} else {
		state.Tab.PendingTx = nil
		state.Tab.ErrorText = err.Error()
	}

	render()
}

type ActCancelTx struct {
}

func (a *ActCancelTx) Run() {
	state.Tab.ProposedTx = nil
	state.Tab.PendingTx = nil

	render()
}

func reloadTxState() {
	tx := state.Tab.PendingTx
	if tx == nil {
		return
	}
	ctx := context.Background()

	receipt, err := client.Ec.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Printf("act reloadTxState error %v", err)
		return
	}
	if receipt == nil {
		return
	}

	log.Printf("act reloadTxState got receipt %s %+v", tx.Hash(), receipt)

	// Transaction confirmed or reverted
	state.Tab.PendingTx = nil
	if receipt.Status == 0 {
		state.Tab.ErrorText = fmt.Sprintf("transaction reverted: %s", tx.Hash())
	}

	render()
}

// Reload context information about the blockchain.
func reloadChainState() {
	state.Chain.Conn = client.ConnStatus()

	render()
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
