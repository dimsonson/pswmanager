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
	LoginPage          string = "LoginPage"
	LoginForm          string = "LoginForm"
	RegisterForm       string = "RegisterForm"
	MainPage           string = "MainPage"
	SelectCreatePage   string = "SelectCreatePage"
	SelectReadPage     string = "SelectReadPage"
	NewTextForm        string = "NewTextForm"
	NewLoginPairForm   string = "NewLoginPairForm"
	NewCardForm        string = "NewCardForm"
	NewBinaryForm      string = "NewBinaryForm"
	ReadTextForm       string = "ReadTextForm"
	ReadLoginPairForm  string = "ReadLoginPairForm"
	ReadCardForm       string = "ReadCardForm"
	ReadBinaryForm     string = "ReadBinaryForm"
	TextSearchResult   string = "TextSearchResult"
	LoginsSearchResult string = "LoginsSearchResult"
	BinarySearchResult string = "BinarySearchResult"
	CardSearchResult   string = "CardSearchResult"

	searchForm string = "searchForm"
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

	flexLoginPairSearchResult *tview.Flex
	flexBinarySearchResult    *tview.Flex
	flexCardSearchResult      *tview.Flex
	flexTextSearchResult      *tview.Flex

	searchForm     *tview.Form
	flexSearchForm *tview.Flex
	searchdata     string
}

type DialogUI struct {
	Confirm *tview.Modal
}

func New() *UI {
	return &UI{}
}

func (ui *UI) NewCustomFlex(primitive tview.Primitive, fixedSize int) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.textMain, 2, 1, false).
		AddItem(primitive, fixedSize, 1, true).
		AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
		AddItem(ui.textMain, 1, 1, false)
	return flex
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
	ui.readLoginPairForm = tview.NewForm()
	ui.readBinaryForm = tview.NewForm()
	ui.readCardForm = tview.NewForm()
	ui.searchForm = tview.NewForm()
	ui.listTextSearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listLoginsSearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listBinarySearchResult = tview.NewList().ShowSecondaryText(false)
	ui.listCardSearchResult = tview.NewList().ShowSecondaryText(false)

	r1 := models.TextRecord{
		Metadata: "test1",
		Text:     "texttext1",
	}
	r2 := models.TextRecord{
		Metadata: "test2",
		Text:     "texttext2",
	}
	ui.TextSearchResult = []models.TextRecord{r1, r2}
	r3 := models.LoginRecord{
		Metadata: "metalogin11",
		Login:    "login11",
	}
	r4 := models.LoginRecord{
		Metadata: "test2",
		Login:    "login23",
	}
	ui.LoginsSearchResult = []models.LoginRecord{r3, r4}
	r5 := models.BinaryRecord{
		Metadata: "bin1",
		Binary:   "binarybody1",
	}
	r6 := models.BinaryRecord{
		Metadata: "test3",
		Binary:   "binarybody33",
	}
	ui.BinarySearchResult = []models.BinaryRecord{r5, r6}
	r7 := models.CardRecord{
		Metadata: "card1",
		Number:   "4410 0545 8945 4589",
	}
	r8 := models.CardRecord{
		Metadata: "card2",
		Number:   "4410 9999 8945 4589",
	}
	ui.CardSearchResult = []models.CardRecord{r7, r8}

	ui.textMain = tview.NewTextView()
	ui.LogWindow = tview.NewTextView()
	ui.TextConfig()
	ui.ListLogin()
	ui.ListMain()
	ui.ListSelectCreate()
	ui.ListSelectRead()
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

	ui.pages.AddPage(ReadTextForm, ui.flexTextRead, true, false)
	ui.pages.AddPage(ReadLoginPairForm, ui.flexLoginPairRead, true, false)
	ui.pages.AddPage(ReadCardForm, ui.flexCardRead, true, false)
	ui.pages.AddPage(ReadBinaryForm, ui.flexBinaryRead, true, false)

	ui.pages.AddPage(TextSearchResult, ui.flexTextSearchResult, true, false)
	ui.pages.AddPage(BinarySearchResult, ui.flexBinarySearchResult, true, false)
	ui.pages.AddPage(CardSearchResult, ui.flexCardSearchResult, true, false)
	ui.pages.AddPage(LoginsSearchResult, ui.flexLoginPairSearchResult, true, false)

	ui.pages.AddPage(searchForm, ui.flexSearchForm, true, false)
}

func (ui *UI) FlexMain() {
	ui.flexMain = ui.NewCustomFlex(ui.listMain, 10)
}

func (ui *UI) ListMain() {
	ui.listMain = tview.NewList().
		AddItem("Create new", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.pages.SwitchToPage(SelectCreatePage)
		}).
		AddItem("Select existed record", "", 'b', func() {
			ui.regform.Clear(true)
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
