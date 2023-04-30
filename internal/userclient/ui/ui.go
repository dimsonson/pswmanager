package ui

import (
	"os"
	"syscall"

	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"

	"github.com/derailed/tcell/v2"
	//"github.com/rivo/tview"
	"github.com/derailed/tview"
)

// var Brand = []string{"MIR", "VISA", "MC", "AMEX"}

const (
	LoginPage        string = "LoginPage"
	LoginForm        string = "LoginForm"
	RegisterForm     string = "RegisterForm"
	MainPage         string = "MainPage"
	SelectPage       string = "SelectPage"
	NewTextForm      string = "NewTextForm"
	NewLoginPairForm string = "NewLoginPairForm"
	NewCardForm      string = "NewCardForm"
	NewBinaryForm    string = "NewBinaryForm"
)

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
	LoginUI
	MainUI
	DialogUI
	CreateUI
	ReadUI
}

type MainUI struct {
	listLogin *tview.List
	listMain  *tview.List
	flexLogin *tview.Flex
	flexMain  *tview.Flex
	textMain  *tview.TextView
	LogWindow *tview.TextView
	pages     *tview.Pages
	MainApp   *tview.Application
}

type LoginUI struct {
	flexLog   *tview.Flex
	flexReg   *tview.Flex
	regform   *tview.Form
	loginform *tview.Form
}

type CreateUI struct {
	flexSelectCreate    *tview.Flex
	listSelectCreate    *tview.List
	flexTextCreate      *tview.Flex
	createTextForm      *tview.Form
	createLoginPairForm *tview.Form
	createBinaryForm    *tview.Form
	createCardForm      *tview.Form
}

type ReadUI struct {
	flexSelectRead         *tview.Flex
	flexMetadataSearchForm *tview.Flex
	listSelect             *tview.List
	createform             *tview.Form
	readform               *tview.Form
}

type DialogUI struct {
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
	ui.createTextForm = tview.NewForm()
	ui.createLoginPairForm = tview.NewForm()
	ui.createBinaryForm = tview.NewForm()
	ui.createCardForm = tview.NewForm()
	ui.readform = tview.NewForm()
	ui.textMain = tview.NewTextView()
	ui.listLogin = tview.NewList()
	ui.listMain = tview.NewList()
	ui.listSelect = tview.NewList()
	ui.LogWindow = tview.NewTextView()
	ui.TextConfig()
	ui.ListLogin()
	ui.ListMain()
	ui.ListSelectCreate()
	ui.FlexLogin()
	ui.FlexMain()
	ui.FlexCreateRead()
	ui.PagesConfig()
}

func (ui *UI) UIRun() {
	if err := ui.MainApp.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) PagesConfig() {
	ui.pages.AddPage(LoginPage, ui.flexLogin, true, true)
	ui.pages.AddPage(LoginForm, ui.flexLog, true, false)
	ui.pages.AddPage(RegisterForm, ui.flexReg, true, false)
	ui.pages.AddPage(MainPage, ui.flexMain, true, false)
	ui.pages.AddPage(SelectPage, ui.flexSelectCreate, true, false)
	ui.pages.AddPage(NewTextForm, ui.flexTextCreate, true, false)
	ui.pages.AddPage(NewLoginPairForm, ui.flexSelectCreate, true, false)
	ui.pages.AddPage(NewCardForm, ui.flexSelectRead, true, false)
	ui.pages.AddPage(NewBinaryForm, ui.flexSelectRead, true, false)
}

func (ui *UI) FlexMain() {
	ui.flexMain = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listMain, 10, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)
	ui.flexMain.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			ui.MainApp.Stop()
		} else if event.Rune() == '1' {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage(MainPage)
		}
		return event
	})
}

func (ui *UI) ListMain() {
	ui.listMain.
		AddItem("Create new", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.pages.SwitchToPage(SelectPage)
		}).
		AddItem("Select existed record", "", 'b', func() {
			ui.regform.Clear(true)
			ui.ListSelectRead()
			ui.pages.SwitchToPage(SelectPage)
		}).
		AddItem("Quit", "", 'q', func() {
			log.Logg = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			ui.MainApp.Stop()
			err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			if err != nil {
				log.Print("stop programm error")
				return
			}
		})
	ui.listMain.SetBorder(true)
	ui.listMain.SetTitle("Main menu")
	ui.listMain.SetTitleAlign(tview.AlignLeft)
	ui.listMain.SetWrapAround(true)
	ui.listMain.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listMain)
}

func (ui *UI) TextConfig() {
	ui.textMain.
		SetTextColor(tcell.ColorMediumBlue).
		SetText(" Password & Secrets Manager 2023 *** press (q) to quit")
	ui.textMain.SetBackgroundColor(tcell.Color102)

	ui.LogWindow.SetScrollable(false)
	ui.LogWindow.SetBackgroundColor(tcell.Color102.TrueColor())
}
