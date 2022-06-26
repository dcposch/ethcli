package action

import (
	"log"
	"strings"

	"dcposch.eth/cli/v2/eth"
	"github.com/ethereum/go-ethereum/common"
)

var (
	client   *eth.Client
	tab      TabState
	renderer func(*TabState)
)

func Init(_client *eth.Client, _renderer func(*TabState)) {
	client = _client
	renderer = _renderer
}

func SetUrl(url string) {
	log.Printf("action SetUrl %s\n", url)

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
	renderer(&tab)
}
