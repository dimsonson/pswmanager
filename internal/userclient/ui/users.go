package ui

import (
	"github.com/derailed/tview"
	"github.com/rs/zerolog/log"
)


func (ui *UI) FlexUsers() {
	ui.flexMain = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.list, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)

	ui.flexLogin = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.loginform, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)

	ui.flexReg = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.text, 2, 1, false).
			AddItem(ui.regform, 10, 1, true).
			AddItem(ui.TextView.LogWindow.SetChangedFunc(func() { ui.App.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.text, 1, 1, false), 0, 2, true)
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
