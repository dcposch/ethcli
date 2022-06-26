package eth

import (
	"fmt"

	"dcposch.eth/cli/v2/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

// A caching Ethereum client. Forwards requests to a JSON RPC client.
type Client struct {
	Ec *ethclient.Client
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
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addr)
	}
	return
}
