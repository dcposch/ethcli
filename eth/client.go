package eth

import (
	"dcposch.eth/cli/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

// A caching Ethereum client. Forwards requests to a JSON RPC client.
type Client struct {
	Ec *ethclient.Client
}

// An address plus ENS name
type NamedAddr struct {
	addr common.Address
	name string
	err  string
}

// Displays an address, eg "0x123456..." or "vitalik.eth" or "⚠️ invalid.eth"
func (a *NamedAddr) Disp() string {
	ret := ""
	if a.err != "" {
		ret += "⚠️ "
	}
	if a.name != "" {
		ret += a.name
	} else {
		hex := a.addr.Hex()
		ret += hex[0:8] + "…" + hex[36:]
	}
	return ret
}

func CreateClient(ethRpcUrl string) *Client {
	ec, err := ethclient.Dial(ethRpcUrl)
	util.Must(err)

	return &Client{
		Ec: ec,
	}
}

func (c *Client) Resolve(ensName string) (addr common.Address, err error) {
	addr, err = ens.Resolve(c.Ec, ensName)
	return
}
