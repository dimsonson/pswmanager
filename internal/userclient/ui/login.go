package ui

import (
	"os"
	"syscall"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"
)

func (ui *UI) FlexLogin() {
	ui.flexLogin = ui.NewCustomFlex(ui.listLogin, 10)
	ui.flexLog = ui.NewCustomFlex(ui.loginform, 10)
	ui.flexReg = ui.NewCustomFlex(ui.regform, 10)
}

func (ui *UI) ListLogin() {
	ui.listLogin = tview.NewList().
		AddItem("Login", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.loginFrm()
			ui.pages.SwitchToPage(LoginForm)
		}).
		AddItem("Registration", "", 'b', func() {
			ui.regform.Clear(true)
			ui.registerFrm()
			ui.pages.SwitchToPage(RegisterForm)
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
	ui.listLogin.SetBorder(true)
	ui.listLogin.SetTitle("Login menu")
	ui.listLogin.SetTitleAlign(tview.AlignLeft)
	ui.listLogin.SetWrapAround(true)
	ui.listLogin.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listLogin)
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
					ui.loginform.SetFocus(0)
					ui.pages.ShowPage(LoginForm)
				},
				func() {
					ui.pages.SwitchToPage(LoginPage)
				})
		}
		if loginpsw.uLogin == "0" {
			log.Print("user login 0")
			ui.ShowOk("Login successful", func() {
				ui.pages.SwitchToPage(MainPage)
			})
		}

		if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
			ui.pages.SwitchToPage(LoginForm)
		}

	})
	ui.loginform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(LoginPage)
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
		if loginpsw.uLogin == "0" {
			log.Print("user reg 0")
			ui.ShowOk("Registration successful. \n Please, keep your credentials safe.", func() {
				ui.pages.SwitchToPage(LoginForm)
			})
		}
		if loginpsw.uLogin == "1" {
			ui.ShowConfirm("Username already exist", "Do you like try again?",
				func() {
					log.Print("user reg 1")
					ui.regform.SetFocus(0)
					ui.pages.ShowPage(RegisterForm)
				},
				func() {
					ui.pages.SwitchToPage(LoginPage)
				})
		}

		if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
			ui.pages.SwitchToPage(LoginForm)
		}
	})
	ui.regform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(LoginPage)
	})
	return ui.regform
}


