package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"os"

	"dcposch.eth/cli/eth"
	"dcposch.eth/cli/util"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Opts struct {
	ethRpcUrl  string
	privateKey *ecdsa.PrivateKey
	logFile    string
}

func main() {
	opts := parseArgsOrExit()
	startLogging(opts.logFile)

	// Connect to Ethereum
	client := eth.CreateClient(opts.ethRpcUrl)

	// Initialize browser state. One-way data flow: action > state > render.
	// act.Init(client, ui.Render)

	vdom, err := client.FrontendRender(common.Address{}, common.HexToAddress("0xa51e457bab8f571ccf72327be0787f32491deac0"), []byte{})
	util.Must(err)

	for _, v := range vdom {
		switch v.TypeHash {
		case eth.TypeText:
			fmt.Printf("Text: %s\n", string(v.Data))
		case eth.TypeInAmount:
			n := abi.ReadInteger(abi.Type{T: abi.UintTy, Size: 64}, v.Data).(uint64)
			fmt.Printf("Amount: %d %x\n", n, v.Data)
		case eth.TypeInDropdown:
			fmt.Printf("Dropdown: %d bytes\n", len(v.Data))
		case eth.TypeInTextbox:
			fmt.Printf("Textbox: %d bytes\n", len(v.Data))
		case eth.TypeButton:
			fmt.Printf("Button: %x\n", v.Data)
		default:
			fmt.Printf("UKNOWN %d: %d bytes\n", v.TypeHash, len(v.Data))
		}
	}

	// Show a terminal dapp browser
	// ui.StartRenderer()
}

// Returns either valid options or exits printing an error message.
func parseArgsOrExit() (r Opts) {
	flag.StringVar(&r.ethRpcUrl, "rpc-url", os.Getenv("ETH_RPC_URL"), "[env ETH_RPC_URL]")
	var privateKeyHex string
	flag.StringVar(&privateKeyHex, "private-key", "", "Account private key")
	flag.StringVar(&r.logFile, "log-file", "", "Debug log file. Default: new temp file.")
	flag.Parse()

	if r.ethRpcUrl == "" {
		flag.Usage()
		fmt.Println("Missing RPC URL")
		os.Exit(2)
	}

	if privateKeyHex != "" {
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
		util.Must(err)
		r.privateKey = privateKey
	}

	return
}

// Log to a temp file. We're about to start tview and cannot log to terminal.
func startLogging(path string) {
	var logFile *os.File
	var err error
	if path == "" {
		logFile, err = os.CreateTemp("", "eth-*")
	} else {
		logFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	}
	util.Must(err)

	log.SetPrefix("ethcli")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds)
	log.Printf("Writing log output to %s", logFile.Name())
	log.SetOutput(logFile)
	log.Println("Hello world")
}
