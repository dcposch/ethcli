package eth

import (
	"time"

	"dcposch.eth/cli/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
	"golang.org/x/net/context"
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

func (c *Client) ConnStatus() ConnStatus {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	cid, err := c.Ec.ChainID(ctx)
	if err != nil {
		return ConnStatus{0, err.Error()}
	} else {
		return ConnStatus{cid.Int64(), ""}
	}
}

type ConnStatus struct {
	ChainID   int64 // TODO: chain name
	ErrorText string
}

func (c *Client) Resolve(ensName string) (addr common.Address, err error) {
	addr, err = ens.Resolve(c.Ec, ensName)
	return
}

// An address plus ENS name
type NamedAddr struct {
	Add  common.Address
	Name string
	Err  string
}

// Displays an address, eg "0x123456..." or "vitalik.eth" or "⚠️ invalid.eth"
func (a *NamedAddr) Disp() string {
	ret := ""
	if a.Err != "" {
		ret += "⚠️ "
	}
	if a.Name != "" {
		ret += a.Name
	} else {
		hex := a.Add.Hex()
		ret += hex[0:8] + "…" + hex[36:]
	}
	return ret
}
