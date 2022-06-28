package ui

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"

	"dcposch.eth/cli/act"
	"dcposch.eth/cli/eth"
	"dcposch.eth/cli/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	fgGreen  = tcell.ColorGreen
	bgDark   = tcell.ColorDarkBlue
	bgGray   = tcell.ColorDarkGray
	bgErr    = tcell.ColorDarkRed
	colReset = tcell.ColorReset
)

var (
	lastState        *act.State
	lastStateStr     string
	app              *tview.Application
	urlInput         *tview.InputField
	chainStatus      *tview.TextView
	mainContent      *tview.Flex
	footerConnStatus *tview.TextView
	footerMain       *tview.TextView
	pages            *tview.Pages
	modalConfirm     *tview.Modal
)

func StartRenderer() {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(32, 80, 0).
		SetBorders(true)

	// Header row
	appLabel := tview.NewTextView().SetTextColor(fgGreen).SetText("ETHEREUM EXPLORER")
	urlInput = tview.NewInputField().SetLabel("ENS or address: ").SetDoneFunc(onDoneUrlInput)
	grid.AddItem(appLabel, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(urlInput, 0, 1, 1, 1, 0, 0, true)
	grid.AddItem(tview.NewTextView(), 0, 2, 1, 1, 0, 0, false)

	// Main row
	chainStatus = tview.NewTextView().SetText("ACCOUNT")
	chainStatus.SetBorderPadding(1, 1, 0, 0)
	mainContent = tview.NewFlex().SetDirection(tview.FlexColumnCSS)
	mainContent.SetBorderPadding(1, 1, 1, 1)
	grid.AddItem(chainStatus, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(mainContent, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 1, 2, 1, 1, 0, 0, false)

	// Footer row
	footerConnStatus = tview.NewTextView()
	footerMain = tview.NewTextView()
	grid.AddItem(footerConnStatus, 2, 0, 1, 1, 0, 0, false)
	grid.AddItem(footerMain, 2, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 2, 2, 1, 1, 0, 0, false)

	modalConfirm = tview.NewModal().
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				txModalConfirm()
			} else {
				txModalCancel()
			}
		})

	pages = tview.NewPages().
		AddPage("main", grid, true, true).
		AddPage("modal", modalConfirm, true, false)

	app = tview.NewApplication().
		SetRoot(pages, true).
		EnableMouse(true)

	// Tab order
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			moveFocus(1)
		} else if event.Key() == tcell.KeyBacktab {
			moveFocus(-1)
		}
		return event
	})

	util.Must(app.Run())
}

func onDoneUrlInput(key tcell.Key) {
	if key == tcell.KeyEnter {
		act.Dispatch(&act.ActSetUrl{Url: urlInput.GetText()})
	} else {
		urlInput.SetText(lastState.Tab.EnteredAddr)
	}
}

var isRendering = false

func Render(state *act.State) {
	// TODO: better diff
	stateStr := fmt.Sprintf("%#v", state)
	if stateStr == lastStateStr {
		return
	}

	app.QueueUpdateDraw(func() {
		isRendering = true
		log.Printf("ui Render %#v URL %s %s err '%s' elems %d", state.Chain,
			state.Tab.EnteredAddr, state.Tab.ContractAddr, state.Tab.ErrorText,
			len(state.Tab.Vdom))

		renderChain(&state.Chain)
		renderTab(&state.Tab)
		renderModal(state)

		lastState = state
		lastStateStr = stateStr
		isRendering = false
	})
}

func txModalConfirm() {
	act.Dispatch(&act.ActExecTx{})
}

func txModalCancel() {
	act.Dispatch(&act.ActCancelTx{})
}

func renderModal(state *act.State) {
	propTx := state.Tab.ProposedTx
	pendTx := state.Tab.PendingTx

	if propTx == nil && pendTx == nil {
		pages.HidePage("modal")
	} else if state.Chain.PrivateKey == nil {
		pages.ShowPage("modal")
		modalConfirm.SetText("You must be logged in to submit transactions.")
	} else {
		pages.ShowPage("modal")
		if propTx != nil {
			modalConfirm.SetText(fmt.Sprintf("Confirm transaction to %s?", propTx.To))
		} else {
			modalConfirm.SetText(fmt.Sprintf("Transaction %s pending...", pendTx.Hash()))
		}
	}
}

func renderTab(tab *act.TabState) {
	footerMain.SetBackgroundColor(bgGray)
	if tab.EnteredAddr == "" {
		footerMain.SetText("Enter a contract address to begin")
	} else if tab.ContractAddr == nil && tab.ErrorText == "" {
		footerMain.SetText("Resolving...")
	} else if tab.ContractAddr != nil {
		footerMain.SetText(fmt.Sprintf("Resolved %s", tab.ContractAddr))
	} else {
		footerMain.SetText(fmt.Sprintf("Error: %s", tab.ErrorText))
		footerMain.SetBackgroundColor(bgErr)
	}

	mainContent.Clear()
	errText := tab.ErrorText
	if errText == "" && tab.Vdom != nil {
		for _, v := range tab.Vdom {
			key := v.DataElem.GetKey()
			inputVal := tab.Inputs[key]
			item, err := createItem(v.DataElem, inputVal)
			if err != nil {
				errText = err.Error()
				break
			}
			mainContent.AddItem(item, 3, 0, false)
		}
	}
	if errText != "" {
		mainContent.Clear()
		errItem := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(tab.ErrorText).
			SetBackgroundColor(bgErr)
		mainContent.AddItem(errItem, 1, 0, false)
	}
	mainContent.AddItem(tview.NewTextView(), 0, 1, false)
}

func resetFocus() {
	app.SetFocus(urlInput)
}

func moveFocus(dir int) {
	log.Printf("moving focus %d", dir)
	focusIx := getFocusIx()
	newIx := focusIx + dir
	if newIx < 0 {
		app.SetFocus(urlInput)
	} else if newIx < mainContent.GetItemCount() {
		app.SetFocus(mainContent.GetItem(newIx))
	}
}

func getFocusIx() int {
	focusE := app.GetFocus()
	focusIx := -1
	for i := 0; i < mainContent.GetItemCount(); i++ {
		if focusE == mainContent.GetItem(i) {
			focusIx = i
			break
		}
	}
	return focusIx
}

func createItem(elem eth.KeyElem, inputVal []byte) (tview.Primitive, error) {
	switch e := elem.(type) {
	case *eth.ElemText:
		return tview.NewTextView().SetText(e.Text), nil
	case *eth.ElemAmount:
		label := padRight(e.Label, 24)

		dec := int(e.Decimals)
		initText := util.ToFixedPrecision(util.DecodeUint(inputVal), dec)
		ret := tview.NewInputField().SetLabel(label).SetText(initText)
		ret.SetDoneFunc(func(key tcell.Key) {
			if isRendering {
				return
			}
			text := ret.GetText()
			log.Printf("amount Done: %d %s", e.Key, text)
			fVal, err := strconv.ParseFloat(text, 64)
			if err != nil {
				ret.SetFieldBackgroundColor(bgErr)
			} else {
				ret.SetFieldBackgroundColor(colReset)
				fVal = math.Round(fVal * math.Pow10(dec))
				val := util.EncodeUint(big.NewInt(int64(fVal)))
				setInput(e.Key, val)
			}
			if key == tcell.KeyEnter {
				moveFocus(1)
			}
		})
		return ret, nil
	case *eth.ElemDropdown:
		label := padRight(e.Label, 24)
		ret := tview.NewDropDown().SetLabel(label)
		initV := util.DecodeUint(inputVal)
		selIx := -1
		for i, opt := range e.Options {
			val := opt.Val
			ret.AddOption(opt.Text, func() {
				if isRendering {
					return
				}
				setInput(e.Key, util.EncodeUint(val))
			})
			if opt.Val.Cmp(initV) == 0 {
				selIx = i
			}
		}
		ret.SetCurrentOption(selIx)
		ret.SetFieldBackgroundColor(bgGray)
		return ret, nil
	case *eth.ElemButton:
		return tview.NewButton(e.Text).SetSelectedFunc(func() {
			if isRendering {
				return
			}
			submit(e.Key)
		}), nil
	default:
		return nil, fmt.Errorf("unimplemented: %t", elem)
	}
}

func setInput(key uint8, val []byte) {
	act.Dispatch(&act.ActSetInput{Key: key, Val: val})
}

func submit(buttonKey uint8) {
	resetFocus()
	act.Dispatch(&act.ActSubmit{ButtonKey: buttonKey})
}

func padRight(label string, width int) string {
	if len(label) > width {
		return label[:width-1] + "…"
	}
	return label + strings.Repeat(" ", width-len(label))
}

func renderChain(chain *act.ChainState) {
	if chain.PrivateKey == nil {
		chainStatus.SetText("Not logged in")
	} else {
		chainStatus.SetText("🔑 " + chain.Account.Disp())
	}

	if chain.Conn.ErrorText == "" {
		statusText := fmt.Sprintf("CONNECTED - %s", strings.ToUpper(chain.Conn.ChainName))
		footerConnStatus.SetText(statusText).SetBackgroundColor(bgDark)
	} else {
		footerConnStatus.SetText("DISCONNECTED").SetBackgroundColor(bgErr)
	}
}
