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

const (
	Menu       string = "Menu"
	Login      string = "Login"
	Register   string = "Register"
	MenuSecond string = "MenuSecond"
	Create     string = "Create"
	Read       string = "Read"
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
	UserUI
	MainUI
	DialogUI
	CreateReadUI
}

type MainUI struct {
	listMain   *tview.List
	listSecond *tview.List
	flexMain   *tview.Flex
	flexSecond *tview.Flex
	textMain   *tview.TextView
	LogWindow  *tview.TextView
	pages      *tview.Pages
	MainApp    *tview.Application
}

type UserUI struct {
	flexLogin *tview.Flex
	flexReg   *tview.Flex
	regform   *tview.Form
	loginform *tview.Form
}

type CreateReadUI struct {
	flexCreate *tview.Flex
	flexRead   *tview.Flex
	createform *tview.Form
	readform   *tview.Form
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
	ui.createform = tview.NewForm()
	ui.readform = tview.NewForm()
	ui.textMain = tview.NewTextView()
	ui.listMain = tview.NewList()
	ui.listSecond = tview.NewList()
	ui.LogWindow = tview.NewTextView()
	ui.TextConfig()
	ui.ListMain()
	ui.ListSecond()
	ui.FlexMain()
	ui.FlexSecond()
	ui.FlexUsers()
	ui.FlexCreateRead()
	ui.PagesConfig()
}

func (ui *UI) UIRun() {
	if err := ui.MainApp.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) PagesConfig() {
	ui.pages.AddPage(Menu, ui.flexMain, true, true)
	ui.pages.AddPage(Login, ui.flexLogin, true, false)
	ui.pages.AddPage(Register, ui.flexReg, true, false)
	ui.pages.AddPage(MenuSecond, ui.flexSecond, true, false)
	ui.pages.AddPage(Create, ui.flexCreate, true, false)
	ui.pages.AddPage(Read, ui.flexRead, true, false)
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
			ui.pages.SwitchToPage(Menu)
		}
		return event
	})
}

func (ui *UI) FlexSecond() {
	ui.flexSecond = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listSecond, 10, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)
	ui.flexSecond.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			ui.MainApp.Stop()
		} else if event.Rune() == '1' {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage(Menu)
		}
		return event
	})
}

func (ui *UI) ListMain() {
	ui.listMain.
		AddItem("Login", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage("Login")
		}).
		AddItem("Registration", "", 'b', func() {
			ui.regform.Clear(true)
			ui.registerFrm()
			ui.pages.SwitchToPage(Register)
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
	ui.listMain.SetBorder(true)
	ui.listMain.SetTitle("Main menu")
	ui.listMain.SetTitleAlign(tview.AlignLeft)
	ui.listMain.SetWrapAround(true)
	ui.listMain.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listMain)
}

func (ui *UI) ListSecond() {
	ui.listSecond.
		AddItem("Create", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.createFrm()
			ui.pages.SwitchToPage(Create)
		}).
		AddItem("Read", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(Read)
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
	ui.listSecond.SetBorder(true)
	ui.listSecond.SetTitle("Main menu")
	ui.listSecond.SetTitleAlign(tview.AlignLeft)
	ui.listSecond.SetWrapAround(true)
	ui.listSecond.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSecond)
}

func (ui *UI) TextConfig() {
	ui.textMain.
		SetTextColor(tcell.ColorMediumBlue).
		SetText(" Password & Secrets Manager 2023 *** press (q) to quit")
	ui.textMain.SetBackgroundColor(tcell.Color102)

	ui.LogWindow.SetScrollable(false)
	ui.LogWindow.SetBackgroundColor(tcell.Color102.TrueColor())
}
