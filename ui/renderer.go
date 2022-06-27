package ui

import (
	"fmt"
	"log"

	"dcposch.eth/cli/act"
	"dcposch.eth/cli/eth"
	"dcposch.eth/cli/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	fgGreen    = tcell.ColorGreen
	bgDarkGray = tcell.ColorDarkGray
	bgErr      = tcell.ColorDarkRed
	colReset   = tcell.ColorReset
)

var (
	lastState       *act.State
	lastStateStr    string
	app             *tview.Application
	urlInput        *tview.InputField
	chainAccount    *tview.TextView
	chainConnStatus *tview.TextView
	mainContent     *tview.TextView
	footer          *tview.TextView
)

func StartRenderer() {
	appLabel := tview.NewTextView().SetTextColor(fgGreen).SetText("ETHEREUM")
	urlInput = tview.NewInputField().SetLabel("ENS or address: ").SetDoneFunc(onDoneUrlInput)

	chainAccount = tview.NewTextView().SetText("ACCOUNT")
	chainConnStatus = tview.NewTextView().SetText("CONN")
	chainStatus := tview.NewFlex().
		SetDirection(tview.FlexColumnCSS).
		AddItem(chainAccount, 0, 1, false).
		AddItem(chainConnStatus, 1, 0, false)

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

	app = tview.NewApplication().
		SetRoot(grid, true).
		EnableMouse(true)

	util.Must(app.Run())
}

func onDoneUrlInput(key tcell.Key) {
	if key == tcell.KeyEnter {
		act.Dispatch(&act.ActSetUrl{Url: urlInput.GetText()})
	} else {
		urlInput.SetText(lastState.Tab.EnteredAddr)
	}
}

func Render(state *act.State) {
	// TODO: better diff
	stateStr := fmt.Sprintf("%#v", state)
	if stateStr == lastStateStr {
		return
	}
	log.Printf("ui Render %#v", state)

	renderTab(&state.Tab)
	renderChain(&state.Chain)
	app.Draw()

	lastState = state
	lastStateStr = stateStr
}

func renderTab(tab *act.TabState) {
	if tab.EnteredAddr == "" {
		footer.SetText("Enter a contract address to begin")
	} else if tab.ContractAddr == nil && tab.ErrorText == "" {
		footer.SetText("Resolving...")
	} else if tab.ContractAddr != nil {
		footer.SetText(fmt.Sprintf("Resolved %s", tab.ContractAddr))
	} else {
		footer.SetText(fmt.Sprintf("Error: %s", tab.ErrorText))
	}

	if tab.Vdom == nil {
		mainContent.SetText(tab.ErrorText)
	} else {
		mainContent.SetText(string(tab.Vdom))
	}
}

func renderChain(chain *act.ChainState) {
	if eth.IsZeroAddr(chain.Account.Addr) {
		chainAccount.SetText("Not logged in")
	} else {
		chainAccount.SetText(chain.Account.Disp())
	}

	if chain.Conn.ErrorText == "" {
		statusText := fmt.Sprintf("CONNECTED - %s", chain.Conn.ChainName)
		chainConnStatus.SetText(statusText).SetBackgroundColor(colReset)
	} else {
		chainConnStatus.SetText("DISCONNECTED").SetBackgroundColor(bgErr)
	}
}
