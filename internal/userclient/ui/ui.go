package ui

import (
	"os"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/derailed/tcell/v2"
	//"github.com/rivo/tview"
	"github.com/derailed/tview"
)

// var Brand = []string{"MIR", "VISA", "MC", "AMEX"}

type (
	confirmFunc func()
	cancelFunc  func()
)

type ULogin struct {
	uLogin string
	uPsw   string
	uid    string
}

type UI struct {
	Lists
	Flexs
	Forms
	Texts
	App
	Pages
	TextView
	ModalForm
}

type Lists struct {
	list *tview.List
}

type Flexs struct {
	flexMain  *tview.Flex
	flexLogin *tview.Flex
	flexReg   *tview.Flex
}

type Forms struct {
	regform   *tview.Form
	loginform *tview.Form
}

type Texts struct {
	text *tview.TextView
}

type App struct {
	MainApp *tview.Application
}

type Pages struct {
	pages *tview.Pages
}

type TextView struct {
	LogWindow *tview.TextView
}

type ModalForm struct {
	Confirm *tview.Modal
}

func New() *UI {
	return &UI{}
}

func (ui *UI) Init() {
	ui.MainApp = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.loginform = tview.NewForm()
	ui.regform = tview.NewForm()
	ui.text = tview.NewTextView()
	ui.list = tview.NewList()
	ui.LogWindow = tview.NewTextView()
	ui.TextConfig()
	ui.ListConfig()
	ui.FlexMain()
	ui.FlexUsers()
	ui.PagesConfig()
}

func (ui *UI) TextConfig() {
	ui.text.
		SetTextColor(tcell.ColorMediumBlue).
		SetText(" Password & Secrets Manager 2023 *** press (q) to quit")
	ui.text.SetBackgroundColor(tcell.Color102)

	ui.LogWindow.SetScrollable(false)
	ui.LogWindow.SetBackgroundColor(tcell.Color102.TrueColor())
}

func (ui *UI) ListConfig() {
	ui.list.
		AddItem("Login", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage("Login")
		}).
		AddItem("Registration", "", 'b', func() {
			ui.regform.Clear(true)
			ui.registerFrm()
			ui.pages.SwitchToPage("Register")
		}).
		AddItem("Quit", "", 'q', func() {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			ui.MainApp.Stop()
			err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			if err != nil {
				log.Print("stop programm error")
				return
			}
		})
	ui.list.SetBorder(true)
	ui.list.SetTitle("Main menu")
	ui.list.SetTitleAlign(tview.AlignLeft)
	ui.list.SetWrapAround(true)
	ui.list.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.list)
}

func (ui *UI) FlexMain() {
	ui.flexMain = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.list, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)
	ui.flexMain.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			ui.MainApp.Stop()
		} else if event.Rune() == '1' {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage("Menu")
		}
		return event
	})
}

func (ui *UI) PagesConfig() {
	ui.pages.AddPage("Menu", ui.flexMain, true, true)
	ui.pages.AddPage("Login", ui.flexLogin, true, false)
	ui.pages.AddPage("Register", ui.flexReg, true, false)
}

func (ui *UI) UIRun() {
	if err := ui.App.MainApp.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

