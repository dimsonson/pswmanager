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
	App *tview.Application
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

func (menu *UI) Init() {
	menu.App.App = tview.NewApplication()
	menu.pages = tview.NewPages()
	menu.loginform = tview.NewForm()
	menu.regform = tview.NewForm()
	menu.text = tview.NewTextView()
	menu.list = tview.NewList()
	menu.TextView.LogWindow = tview.NewTextView()

}

func (ui *UI) TextConfig() {
	ui.text.
		SetTextColor(tcell.ColorMediumBlue).
		SetText(" Password & Secrets Manager 2023 *** press (q) to quit")
	ui.text.SetBackgroundColor(tcell.Color102)

	ui.TextView.LogWindow.SetScrollable(false)
	ui.TextView.LogWindow.SetBackgroundColor(tcell.Color102.TrueColor())
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
			ui.App.App.Stop()
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
	ui.App.App.SetFocus(ui.list)
}

func (ui *UI) FlexConfig() {
	ui.flexMain = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.list, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.App.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)

	ui.flexLogin = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.loginform, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.App.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)

	ui.flexReg = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.regform, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.App.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)

	ui.flexMain.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			ui.App.App.Stop()
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

func (ui *UI) loginFrm() *tview.Form {
	loginpsw := ULogin{}
	ui.loginform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		loginpsw.uLogin = ulogin
	})
	ui.loginform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		loginpsw.uPsw = upsw
	})
	ui.loginform.AddButton("Login", func() {
		if loginpsw.uLogin == "1" {
			ui.ShowConfirm("Wrong password or username", "Do you like try again?",
				func() {
					log.Print("user login 1")

					//ui.loginform.Clear(true)
					//ui.loginFrm()
					ui.loginform.SetFocus(0) //.SetFocus(0)
					ui.pages.ShowPage("Login")
				},
				func() {
					ui.pages.SwitchToPage("Menu")
				})
		}
		if loginpsw.uLogin == "0" {
			log.Print("user login 0")
			ui.ShowOk("Login successful", func() {
				ui.pages.SwitchToPage("Menu")
			})
			//ui.loginform.Clear(true)
			//ui.pages.SwitchToPage("Menu")
		}

		if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
			ui.pages.SwitchToPage("Menu")
		}
	})
	ui.loginform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage("Menu")
	})
	return ui.loginform
}

func (ui *UI) registerFrm() *tview.Form {
	loginpsw := ULogin{}
	ui.regform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		loginpsw.uLogin = ulogin
	})
	ui.regform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		loginpsw.uPsw = upsw
	})
	ui.regform.AddButton("Register", func() {
		ui.pages.SwitchToPage("Menu")
	})
	ui.regform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage("Menu")
	})
	return ui.regform
}

func (ui *UI) UIRun() {
	if err := ui.App.App.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) ShowConfirm(title, msg string, ack confirmFunc, cancel cancelFunc) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.Color110).
		SetButtonTextColor(tcell.Color111).
		SetLabelColor(tcell.Color112.TrueColor()).
		SetFieldTextColor(tcell.Color114.TrueColor())
	f.AddButton("Cancel", func() {
		ui.dismiss(ui.Pages)
		cancel()
	})
	f.AddButton("OK", func() {
		ui.dismiss(ui.Pages)
		ack()
	})
	for i := 0; i < 2; i++ {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColorActivated(tcell.Color116.TrueColor())
		b.SetLabelColorActivated(tcell.Color117.TrueColor())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm(title, f)
	modal.SetText(msg)
	modal.SetTextColor(tcell.Color119.TrueColor())
	modal.SetDoneFunc(func(int, string) {
		ui.dismiss(ui.Pages)
		cancel()
	})
	ui.pages.AddPage("dialogKey", modal, false, false)
	ui.pages.ShowPage("dialogKey")
}

func (ui *UI) dismiss(pages Pages) {
	ui.pages.RemovePage("dialogKey")
}

func (ui *UI) ShowOk(msg string, ok confirmFunc) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.Color110).
		SetButtonTextColor(tcell.Color111).
		SetLabelColor(tcell.Color112.TrueColor()).
		SetFieldTextColor(tcell.Color114.TrueColor())
	f.AddButton("OK", func() {
		ui.dismiss(ui.Pages)
		ok()
	})
	if b := f.GetButton(0); b != nil {
		b.SetBackgroundColorActivated(tcell.Color116.TrueColor())
		b.SetLabelColorActivated(tcell.Color117.TrueColor())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("", f)
	modal.SetText(msg)
	modal.SetTextColor(tcell.ColorOrangeRed)
	modal.SetDoneFunc(func(int, string) {
		ui.dismiss(ui.Pages)
		ok()
	})
	ui.pages.AddPage("dialogKey", modal, false, false)
	ui.pages.ShowPage("dialogKey")
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
