package act

import (
	"strings"

	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum/common"
)

type Action interface {
	Run()
}

type ActSetUrl struct {
	Url string
}

func (a *ActSetUrl) Run() {
	url := a.Url
	tab := &state.Tab

	tab.EnteredAddr = url
	tab.ErrorText = ""
	tab.ContractAddr = nil

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

func reloadChainState() {
	state.Chain.Conn = client.ConnStatus()
	state.Chain.Account = eth.NamedAddr{}
}

func reloadTab() {
	if state.Tab.ContractAddr == nil {
		return
	}

	appState := []byte{}
	vdom, err := client.FrontendRender(state.Chain.Account.Addr, *state.Tab.ContractAddr, appState)
	if err == nil {
		// TODO: vdom diffing
		state.Tab.Vdom = vdom
		state.Tab.ErrorText = ""
	} else {
		state.Tab.Vdom = nil
		state.Tab.ErrorText = err.Error()
	}

	render()
}

func render() {
	renderer(&state)
}
