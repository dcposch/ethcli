package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"dcposch.eth/cli/v2/action"
	"dcposch.eth/cli/v2/eth"
	"dcposch.eth/cli/v2/ui"
	"dcposch.eth/cli/v2/util"
)

type Opts struct {
	ethRpcUrl string
}

func main() {
	opts := parseArgsOrExit()
	startLogging()

	// Connect to Ethereum
	client := eth.CreateClient(opts.ethRpcUrl)

	// Initialize browser state. One-way data flow: action > state > render.
	action.Init(client, ui.Render)

	// Show a terminal dapp browser
	ui.StartRenderer()
}

// Returns either valid options or exits printing an error message.
func parseArgsOrExit() (r Opts) {
	flag.StringVar(&r.ethRpcUrl, "rpc-url", os.Getenv("ETH_RPC_URL"), "[env ETH_RPC_URL]")
	flag.Parse()

	if r.ethRpcUrl == "" {
		flag.Usage()
		fmt.Println("Missing RPC URL")
		os.Exit(2)
	}

	return
}

// Log to a temp file. We're about to start tview and cannot log to terminal.
func startLogging() {
	logFile, err := os.CreateTemp("", "eth-*")
	util.Must(err)
	log.Printf("Writing log output to %s", logFile.Name())
	log.SetOutput(logFile)
}
