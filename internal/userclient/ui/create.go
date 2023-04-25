package ui

import (
	"os"
	"syscall"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (ui *UI) FlexCreateRead() {
	ui.flexCreate = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listSelect, 12, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)

	ui.flexRead = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.readform, 10, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)

	ui.flexSelect = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listSelect, 10, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)
}

func (ui *UI) ListSelect() {
	ui.listSelect.
		AddItem("Text", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.createFrm()
			ui.pages.SwitchToPage(Create)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(Read)
		}).
		AddItem("Binary", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(Read)
		}).
		AddItem("Card", "", 'b', func() {
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
	ui.listSelect.SetBorder(true)
	ui.listSelect.SetTitle("Main menu")
	ui.listSelect.SetTitleAlign(tview.AlignLeft)
	ui.listSelect.SetWrapAround(true)
	ui.listSelect.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSelect)
}

func (ui *UI) createFrm() *tview.Form {
	loginpsw := ULogin{}
	// ui.createform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
	// 	loginpsw.uLogin = ulogin
	// })
	// ui.createform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
	// 	loginpsw.uPsw = upsw
	// })
	ui.createform.AddButton("New Login Item", func() {
		if loginpsw.uLogin == "1" {
			ui.ShowConfirm("Wrong password or username", "Do you like try again?",
				func() {
					log.Print("user login 1")
					ui.createform.SetFocus(0)
					ui.pages.ShowPage(Login)
				},
				func() {

					ui.pages.SwitchToPage(Menu)
				})
		}
		if loginpsw.uLogin == "0" {
			log.Print("user login 0")
			ui.ShowOk("Login successful", func() {
				ui.pages.SwitchToPage(Menu)
			})
		}
		if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
			ui.pages.SwitchToPage(Menu)
		}
	})
	ui.createform.AddButton("New Text Item", func() {
		ui.pages.SwitchToPage(Menu)
	})
	ui.createform.AddButton("New Binary Item", func() {
		ui.pages.SwitchToPage(Menu)
	})
	ui.createform.AddButton("New Card Item", func() {
		ui.pages.SwitchToPage(Menu)
	})
	ui.createform.SetButtonsAlign(0)

	return ui.createform
}

func (ui *UI) readFrm() *tview.Form {
	loginpsw := ULogin{}
	ui.readform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
		loginpsw.uLogin = ulogin
	})
	ui.readform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
		loginpsw.uPsw = upsw
	})
	ui.readform.AddButton("Register", func() {
		if loginpsw.uLogin == "0" {
			log.Print("user reg 0")
			ui.ShowOk("Registration successful. \n Please, keep your credentials safe.", func() {

				ui.pages.SwitchToPage(Menu)
			})
		}
		if loginpsw.uLogin == "1" {
			ui.ShowConfirm("Username already exist", "Do you like try again?",
				func() {
					log.Print("user reg 1")
					ui.readform.SetFocus(0)
					ui.pages.ShowPage(Register)
				},
				func() {
					ui.pages.SwitchToPage(Menu)
				})
		}

		if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
			ui.pages.SwitchToPage(Menu)
		}
	})
	ui.readform.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(Menu)
	})
	return ui.readform
}
