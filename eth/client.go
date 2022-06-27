package eth

import (
	"fmt"
	"log"
	"strings"
	"time"

	"dcposch.eth/cli/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
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
		return ConnStatus{0, "", err.Error()}
	} else {
		name := params.NetworkNames[cid.String()]
		if name == "" {
			name = fmt.Sprintf("CHAIN ID %d", cid)
		}
		return ConnStatus{cid.Int64(), name, ""}
	}
}

type ConnStatus struct {
	ChainID   int64
	ChainName string
	ErrorText string
}

func (c *Client) Resolve(ensName string) (addr common.Address, err error) {
	addr, err = ens.Resolve(c.Ec, ensName)
	return
}

const abiIFrontendJson = `[{"inputs":[{"internalType":"bytes","name":"appState","type":"bytes"},{"components":[{"internalType":"uint256","name":"buttonId","type":"uint256"},{"internalType":"bytes[]","name":"inputs","type":"bytes[]"}],"internalType":"struct Action","name":"action","type":"tuple"}],"name":"act","outputs":[{"internalType":"bytes","name":"newAppState","type":"bytes"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"appState","type":"bytes"}],"name":"render","outputs":[{"components":[{"internalType":"uint64","name":"typeHash","type":"uint64"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"struct VdomElem[]","name":"vdom","type":"tuple[]"}],"stateMutability":"view","type":"function"}]`

var abiIFrontend *abi.ABI

func (c *Client) FrontendRender(fromAddr, contractAddr common.Address, appState []byte) (vdom []VdomElem, err error) {
	if abiIFrontend == nil {
		abiObj, err := abi.JSON(strings.NewReader(abiIFrontendJson))
		util.Must(err)
		abiIFrontend = &abiObj
	}

	data, err := abiIFrontend.Pack("render", appState)
	if err != nil {
		return nil, err
	}

	log.Printf("eth FrontendRender %s", contractAddr)
	callMsg := ethereum.CallMsg{
		From: fromAddr,
		To:   &contractAddr,
		Data: data,
	}
	vdomBytes, err := c.Ec.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	err = abiIFrontend.UnpackIntoInterface(&vdom, "render", vdomBytes)
	if err != nil {
		return nil, err
	}

	for i, v := range vdom {
		switch v.TypeHash {
		case TypeText:
			v.DataStruct = string(v.Data)
		case TypeInAmount:
			var dat DataAmount
			ParseTuple(v.Data, ElemsInAmount, &dat)
			v.DataStruct = dat
		case TypeInDropdown:
			continue
		case TypeInTextbox:
			continue
		case TypeButton:
			var dat DataBtnAction
			ParseTuple(v.Data, ElemsButton, &dat)
		default:
			continue
		}
		vdom[i] = v
	}
	return
}
