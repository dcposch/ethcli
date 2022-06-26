package ui

import (
	"fmt"
	"log"

	"dcposch.eth/cli/act"
	"dcposch.eth/cli/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	fgGreen    = tcell.ColorGreen
	bgDarkGray = tcell.ColorDarkGray
)

var (
	lastState   *act.State
	urlInput    *tview.InputField
	chainStatus *tview.TextView
	mainContent *tview.TextView
	footer      *tview.TextView
)

func StartRenderer() {
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
	log.Printf("ui Render %#v", state)
	renderTab(&state.Tab)
	lastState = state
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
}
