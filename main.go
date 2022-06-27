package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"os"

	"dcposch.eth/cli/act"
	"dcposch.eth/cli/eth"
	"dcposch.eth/cli/ui"
	"dcposch.eth/cli/util"
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
	act.Init(client, ui.Render)

	// Show a terminal dapp browser
	ui.StartRenderer()
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
