package ui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// var Brand = []string{"MIR", "VISA", "MC", "AMEX"}

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
	App *tview.Application
}

type Pages struct {
	pages *tview.Pages
}

type TextView struct {
	LogWindow *tview.TextView
}

func New() *UI {
	return &UI{}
}

func (menu *UI) Init() {
	menu.App.App = tview.NewApplication()
	menu.pages = tview.NewPages()
	menu.loginform = tview.NewForm()
	menu.regform = tview.NewForm()
	menu.text = tview.NewTextView()
	menu.list = tview.NewList()
	menu.TextView.LogWindow = tview.NewTextView()
}

func (mn *UI) TextConfig() {
	mn.text.
		SetTextColor(tcell.ColorMediumBlue).
		SetText(" Password & Secrets Manager 2023 *** press (q) to quit")
	mn.text.SetBackgroundColor(tcell.Color102)

	mn.TextView.LogWindow.SetScrollable(false)
	mn.TextView.LogWindow.SetBackgroundColor(tcell.Color102.TrueColor())
}

func (mn *UI) ListConfig() {
	mn.list.
		AddItem("Login", "", 'a', func() {
			mn.loginform.Clear(true)
			mn.loginFrm()
			mn.pages.SwitchToPage("Login")
		}).
		AddItem("Registration", "", 'b', func() {
			mn.regform.Clear(true)
			mn.registerFrm()
			mn.pages.SwitchToPage("Register")
		}).
		AddItem("Quit", "", 'q', func() {
			os.Interrupt.Signal()
			mn.App.App.Stop()
		})
	mn.list.SetBorder(true)
	mn.list.SetTitle("Main menu")
	mn.list.SetTitleAlign(tview.AlignLeft)
	mn.list.SetWrapAround(true)
	mn.list.SetBackgroundColor(tcell.Color108)
	mn.App.App.SetFocus(mn.list)
}

func (mn *UI) FlexConfig() {
	mn.flexMain = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(mn.text, 2, 1, false).
			AddItem(mn.list, 10, 1, true).
			AddItem(mn.TextView.LogWindow.SetChangedFunc(func() { mn.App.App.Draw() }), 10, 0, false).
			AddItem(mn.text, 1, 1, false), 0, 2, true)

	mn.flexLogin = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(mn.text, 2, 1, false).
			AddItem(mn.loginform, 10, 1, true).
			AddItem(mn.TextView.LogWindow.SetChangedFunc(func() { mn.App.App.Draw() }), 10, 0, false).
			AddItem(mn.text, 1, 1, false), 0, 2, true)

	mn.flexReg = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(mn.text, 2, 1, false).
			AddItem(mn.regform, 10, 1, true).
			AddItem(mn.TextView.LogWindow.SetChangedFunc(func() { mn.App.App.Draw() }), 10, 0, false).
			AddItem(mn.text, 1, 1, false), 0, 2, true)

	mn.flexMain.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			mn.App.App.Stop()
		} else if event.Rune() == '1' {
			mn.loginform.Clear(true)
			mn.loginFrm()
			mn.pages.SwitchToPage("Menu")
		}
		return event
	})
}

func (mn *UI) PagesConfig() {
	mn.pages.AddPage("Menu", mn.flexMain, true, true)
	mn.pages.AddPage("Login", mn.flexLogin, true, false)
	mn.pages.AddPage("Register", mn.flexReg, true, false)
}

func (mn *UI) loginFrm() *tview.Form {
	loginpsw := ULogin{}
	mn.loginform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		loginpsw.uLogin = ulogin
	})
	mn.loginform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		loginpsw.uPsw = upsw
	})
	mn.loginform.AddButton("Login", func() {
		if loginpsw.uLogin == "1" {

		}
		mn.pages.SwitchToPage("Menu")
	})
	mn.loginform.AddButton("Cancel", func() {
		mn.pages.SwitchToPage("Menu")
	})
	return mn.loginform
}

func (mn *UI) registerFrm() *tview.Form {
	loginpsw := ULogin{}
	mn.regform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		loginpsw.uLogin = ulogin
	})
	mn.regform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		loginpsw.uPsw = upsw
	})
	mn.regform.AddButton("Register", func() {
		mn.pages.SwitchToPage("Menu")
	})
	mn.regform.AddButton("Cancel", func() {
		mn.pages.SwitchToPage("Menu")
	})
	return mn.regform
}

func (ui *UI) UIRun() {
	if err := ui.App.App.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// func main() {
// 	ui := New()
// 	ui.Init()
// 	ui.textConfig()
// 	log.SetOutput(ui.textView)
// 	ui.listConfig()
// 	ui.flexConfig()
// 	ui.pagesConfig()

// 	go func() {
// 		for n := 0; n < 100; n++ {
// 			time.Sleep(time.Second * 2)
// 			fmt.Fprintf(ui.textView, "string %d \n", n)
// 		}
// 	}()

// 	if err := ui.App.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
// 		panic(err)
// 	}
// }
