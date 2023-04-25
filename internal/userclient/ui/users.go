package ui

import (
	"github.com/derailed/tview"
	"github.com/rs/zerolog/log"
)

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
