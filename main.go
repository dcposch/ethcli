package main

import (
	"flag"
	"log"
	"os"

	"dcposch.eth/cli/v2/eth"
	"dcposch.eth/cli/v2/ui"
	"dcposch.eth/cli/v2/util"
)

type Opts struct {
	ethRpcUrl string
}

func main() {
	// Handle options
	opts := parseArgs()

	logFile, err := os.CreateTemp("", "eth-*")
	util.Must(err)
	log.Printf("Writing log output to %s", logFile.Name())
	log.SetOutput(logFile)

	// Connect to the chain
	client := eth.CreateClient(opts.ethRpcUrl)

	// Show a browser
	app := ui.CreateBrowserApp(client)
	util.Must(app.Run())
}

func parseArgs() (r Opts) {
	flag.StringVar(&r.ethRpcUrl, "rpc-url", os.Getenv("ETH_RPC_URL"), "Ethereum JSON RPC URL")
	flag.Parse()
	return
}
