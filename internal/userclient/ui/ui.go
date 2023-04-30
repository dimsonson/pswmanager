package ui

import (
	"os"
	"syscall"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
)

const (
	LoginPage         string = "LoginPage"
	LoginForm         string = "LoginForm"
	RegisterForm      string = "RegisterForm"
	MainPage          string = "MainPage"
	SelectCreatePage  string = "SelectCreatePage"
	SelectReadPage    string = "SelectReadPage"
	NewTextForm       string = "NewTextForm"
	NewLoginPairForm  string = "NewLoginPairForm"
	NewCardForm       string = "NewCardForm"
	NewBinaryForm     string = "NewBinaryForm"
	ReadTextForm      string = "ReadTextForm"
	ReadLoginPairForm string = "ReadLoginPairForm"
	ReadCardForm      string = "ReadCardForm"
	ReadBinaryForm    string = "ReadBinaryForm"
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
	listSelectCreate    *tview.List
	flexSelectCreate    *tview.Flex
	flexTextCreate      *tview.Flex
	flexLoginPairCreate *tview.Flex
	flexBinaryCreate    *tview.Flex
	flexCardCreate      *tview.Flex
	createTextForm      *tview.Form
	createLoginPairForm *tview.Form
	createBinaryForm    *tview.Form
	createCardForm      *tview.Form
}

type ReadUI struct {
	flexSelectRead         *tview.Flex
	flexMetadataSearchForm *tview.Flex
	listSelectRead         *tview.List
	flexTextRead           *tview.Flex
	flexLoginPairRead      *tview.Flex
	flexBinaryRead         *tview.Flex
	flexCardRead           *tview.Flex
	readTextForm           *tview.Form
	readLoginPairForm      *tview.Form
	readBinaryForm         *tview.Form
	readCardForm           *tview.Form
	listTextSearchResult   *tview.List
	listLoginsSearchResult *tview.List
	listBinarySearchResult *tview.List
	listCardSearchResult   *tview.List
	TextSearchResult       []models.TextRecord
	LoginsSearchResult     []models.LoginRecord
	BinarySearchResult     []models.BinaryRecord
	CardSearchResult       []models.CardRecord
}

type DialogUI struct {
	Confirm *tview.Modal
}

func New() *UI {
	return &UI{}
}

func (ui *UI) NewAppFlex(primitive tview.Primitive, fixedSize int) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.textMain, 2, 1, false).
		AddItem(primitive, fixedSize, 1, true).
		AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
		AddItem(ui.textMain, 1, 1, false)
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
	ui.readTextForm = tview.NewForm()
	ui.listTextSearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listLoginsSearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listBinarySearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listCardSearchResult = tview.NewList().ShowSecondaryText(false)
	ui.textMain = tview.NewTextView()
	ui.listLogin = tview.NewList()
	ui.listMain = tview.NewList()
	ui.LogWindow = tview.NewTextView()
	ui.TextConfig()
	ui.ListLogin()
	ui.ListMain()
	ui.ListSelectCreate()
	ui.FlexLogin()
	ui.FlexMain()
	ui.FlexCreate()
	ui.FlexRead()
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
	ui.pages.AddPage(SelectCreatePage, ui.flexSelectCreate, true, false)
	ui.pages.AddPage(NewTextForm, ui.flexTextCreate, true, false)
	ui.pages.AddPage(NewLoginPairForm, ui.flexLoginPairCreate, true, false)
	ui.pages.AddPage(NewCardForm, ui.flexCardCreate, true, false)
	ui.pages.AddPage(NewBinaryForm, ui.flexBinaryCreate, true, false)
	ui.pages.AddPage(SelectReadPage, ui.flexSelectRead, true, false)

	// ui.pages.AddPage(ReadTextForm, ui.flexTextRead, true, false)
	// ui.pages.AddPage(ReadLoginPairForm, ui.flexLoginPairRead, true, false)
	// ui.pages.AddPage(ReadCardForm, ui.flexCardRead, true, false)
	// ui.pages.AddPage(ReadBinaryForm, ui.flexBinaryRead, true, false)
}

func (ui *UI) FlexMain() {
	ui.flexMain = ui.NewAppFlex(ui.listMain, 10)
}

func (ui *UI) ListMain() {
	ui.listMain.
		AddItem("Create new", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.pages.SwitchToPage(SelectCreatePage)
		}).
		AddItem("Select existed record", "", 'b', func() {
			ui.regform.Clear(true)
			ui.ListSelectRead()
			ui.pages.SwitchToPage(SelectReadPage)
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
	ui.LogWindow.SetBackgroundColor(tcell.Color102)
}
