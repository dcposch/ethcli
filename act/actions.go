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

	render()
}

func render() {
	renderer(&state)
}

func reloadChainState() {
	state.Chain.Conn = client.ConnStatus()
	state.Chain.Account = eth.NamedAddr{}
}

func reloadTabState() {
}
