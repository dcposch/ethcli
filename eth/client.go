package eth

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"dcposch.eth/cli/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	ens "github.com/wealdtech/go-ens/v3"
	"golang.org/x/net/context"
)

// A caching Ethereum client. Forwards requests to a JSON RPC client.
type Client struct {
	Ec             *ethclient.Client
	LastConnStatus ConnStatus
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
		c.LastConnStatus = ConnStatus{0, "", err.Error()}
	} else {
		name := params.NetworkNames[cid.String()]
		if name == "" {
			name = fmt.Sprintf("CHAIN ID %d", cid)
		}
		c.LastConnStatus = ConnStatus{cid.Int64(), name, ""}
	}
	return c.LastConnStatus
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

const abiIFrontendJson = `[{"inputs":[{"internalType":"bytes","name":"appState","type":"bytes"},{"components":[{"internalType":"uint256","name":"buttonKey","type":"uint256"},{"internalType":"bytes[]","name":"inputs","type":"bytes[]"}],"internalType":"struct Action","name":"action","type":"tuple"}],"name":"act","outputs":[{"internalType":"bytes","name":"newAppState","type":"bytes"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"appState","type":"bytes"}],"name":"render","outputs":[{"components":[{"internalType":"uint64","name":"typeHash","type":"uint64"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"struct VElem[]","name":"vdom","type":"tuple[]"}],"stateMutability":"view","type":"function"}]`

var abiIFrontend = parseAbi(abiIFrontendJson)

func parseAbi(json string) *abi.ABI {
	abiObj, err := abi.JSON(strings.NewReader(abiIFrontendJson))
	util.Must(err)
	return &abiObj
}

func (c *Client) FrontendRender(fromAddr, contractAddr common.Address, appState []byte) (vdom []VElem, err error) {
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
		err := unpackElem(&v)
		if err != nil {
			return nil, err
		}
		vdom[i] = v
	}
	return
}

func (c *Client) FrontendSubmit(fromAddr, contractAddr common.Address, appState []byte, action ButtonAction) (msg *ethereum.CallMsg, err error) {
	abiAction := struct {
		ButtonKey *big.Int
		Inputs    [][]byte
	}{
		ButtonKey: big.NewInt(int64(action.ButtonKey)),
		Inputs:    action.Inputs,
	}
	data, err := abiIFrontend.Pack("act", appState, abiAction)
	if err != nil {
		return nil, err
	}

	log.Printf("eth FrontendSubmit %s", contractAddr)
	callMsg := ethereum.CallMsg{
		From: fromAddr,
		To:   &contractAddr,
		Data: data,
	}
	_, err = c.Ec.CallContract(context.Background(), callMsg, nil)

	return &callMsg, err
}

func (c *Client) Execute(msg *ethereum.CallMsg, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	ctx := context.Background()

	nonce, err := c.Ec.PendingNonceAt(ctx, msg.From)
	if err != nil {
		return nil, fmt.Errorf("nonce %s", err)
	}
	gasPrice, err := c.Ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("price %s", err)
	}
	gas, err := c.Ec.EstimateGas(ctx, *msg)
	if err != nil {
		return nil, fmt.Errorf("gas %s", err)
	}

	// Infura gives "method not suppported"
	// gasTipCap, err := c.Ec.SuggestGasTipCap(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("tip %s", err)
	// }
	gasTipCap := big.NewInt(2_000_000_000) // 2 gwei
	if gasTipCap.Cmp(gasPrice) > 0 {
		gasTipCap.SetInt64(0)
	}

	chainID := big.NewInt(c.LastConnStatus.ChainID)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: gasTipCap,
		Gas:       gas,
		To:        msg.To,
		Value:     msg.Value,
		Data:      msg.Data,
	})

	log.Printf("eth SIGNING TRANSACTION. chain %d nonce %d fee cap %s tip %s gas %d from %s to %s",
		chainID,
		nonce,
		gasPrice,
		gasTipCap,
		gas,
		msg.From,
		msg.To,
	)

	txS, err := types.SignTx(tx, types.NewLondonSigner(chainID), prv)
	if err != nil {
		return nil, err
	}

	return txS, c.Ec.SendTransaction(ctx, txS)
}

func unpackElem(v *VElem) error {
	// slog.Printf("UNPACKING %d: %x", v.TypeHash, v.Data)
	switch v.TypeHash {
	case TypeText:
		v.DataElem = &ElemText{}
		return ParseTuple(v.Data, PropsText, v.DataElem)
	case TypeInAmount:
		v.DataElem = &ElemAmount{}
		return ParseTuple(v.Data, PropsAmount, v.DataElem)
	case TypeInDropdown:
		v.DataElem = &ElemDropdown{}
		return ParseTuple(v.Data, PropsDropdown, v.DataElem)
	case TypeButton:
		v.DataElem = &ElemButton{}
		return ParseTuple(v.Data, PropsButton, v.DataElem)
	case TypeInTextbox:
	default:
		return fmt.Errorf("unsupported elem %d", v.TypeHash)
	}
	return nil
}
