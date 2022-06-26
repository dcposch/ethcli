package ui

import (
	"fmt"
	"log"
	"strings"

	"dcposch.eth/cli/v2/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	fgGreen    = tcell.ColorGreen
	bgDarkGray = tcell.ColorDarkGray
)

var (
	client      *eth.Client
	urlInput    *tview.InputField
	chainStatus *tview.TextView
	mainContent *tview.TextView
	footer      *tview.TextView
	tab         Tab
)

func CreateBrowserApp(_client *eth.Client) *tview.Application {
	client = _client // TODO

	appLabel := tview.NewTextView().SetTextColor(fgGreen).SetText("ETHEREUM")
	urlInput = tview.NewInputField().SetLabel("ENS or address: ").SetDoneFunc(onDoneUrlInput)

	chainStatus = tview.NewTextView()

	mainContent = tview.NewTextView().SetTextAlign(tview.AlignCenter)

	footer = tview.NewTextView()
	footer.SetBackgroundColor(bgDarkGray)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(32, 0).
		SetBorders(true)

	grid.AddItem(appLabel, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(urlInput, 0, 1, 1, 1, 0, 0, true)
	grid.AddItem(chainStatus, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(mainContent, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(footer, 2, 0, 1, 2, 0, 0, false)

	app := tview.NewApplication().SetRoot(grid, true).EnableMouse(true)

	return app
}

func onDoneUrlInput(key tcell.Key) {
	log.Printf("WTF %d", key)
	if key == tcell.KeyEnter {
		actionSetUrl(urlInput.GetText())
	} else {
	}
}

func actionSetUrl(url string) {
	log.Printf("DBG SETTING URL %s\n", url)

	tab.enteredAddr = url
	tab.errorText = ""
	tab.contractAddr = nil

	if strings.HasSuffix(url, ".eth") {
		render()
		result, err := client.Resolve(url)
		if err != nil {
			tab.errorText = err.Error()
		} else {
			tab.contractAddr = &result
		}
	} else if strings.HasPrefix(url, "0x") {
		addr := common.HexToAddress(url)
		tab.contractAddr = &addr
	} else {
		tab.enteredAddr = ""
	}

	render()
}

func render() {
	log.Printf("RENDERING")
	if tab.enteredAddr == "" {
		mainContent.SetText("Enter a contract address to begin")
	} else if tab.contractAddr == nil && tab.errorText == "" {
		mainContent.SetText("Resolving...")
	} else if tab.contractAddr != nil {
		mainContent.SetText(fmt.Sprintf("Resolved %s", tab.contractAddr))
	} else {
		mainContent.SetText(fmt.Sprintf("Error: %s", tab.errorText))
	}
}
