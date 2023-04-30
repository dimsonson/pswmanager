package ui

import (
	"os"
	"syscall"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"
)

func (ui *UI) FlexCreateRead() {
	ui.flexTextCreate = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.createTextForm, 12, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)

	ui.flexSelectCreate = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listSelect, 12, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)

	ui.flexSelectRead = tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ui.textMain, 2, 1, false).
			AddItem(ui.listSelect, 12, 1, true).
			AddItem(ui.LogWindow.SetChangedFunc(func() { ui.MainApp.Draw() }), 10, 0, false).
			AddItem(ui.textMain, 1, 1, false), 0, 2, true)
}

func (ui *UI) ListSelectCreate() {
	ui.listSelect.
		AddItem("Text", "", 'a', func() {
			ui.createTextForm.Clear(true)
			ui.createTextFrm()
			ui.pages.SwitchToPage(NewTextForm)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(NewLoginPairForm)
		}).
		AddItem("Binary", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(NewBinaryForm)
		}).
		AddItem("Card", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(NewCardForm)
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
	ui.listSelect.SetBorder(true)
	ui.listSelect.SetTitle("menu")
	ui.listSelect.SetTitleAlign(tview.AlignLeft)
	ui.listSelect.SetWrapAround(true)
	ui.listSelect.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSelect)
}

func (ui *UI) ListSelectRead() {
	ui.listSelect = tview.NewList()
	ui.listSelect.
		AddItem("Text Item", "", 'a', func() {
			ui.loginform.Clear(true)
			ui.createTextFrm()
			ui.pages.SwitchToPage(LoginPage)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(LoginPage)
		}).
		AddItem("Binary", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(LoginPage)
		}).
		AddItem("Card", "", 'b', func() {
			ui.regform.Clear(true)
			ui.readFrm()
			ui.pages.SwitchToPage(LoginPage)
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
	ui.listSelect.SetBorder(true)
	ui.listSelect.SetTitle("menu")
	ui.listSelect.SetTitleAlign(tview.AlignLeft)
	ui.listSelect.SetWrapAround(true)
	ui.listSelect.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSelect)
}

func (ui *UI) createTextFrm() *tview.Form {
	textRecord := models.TextRecord{}
	ui.createTextForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		textRecord.Metadata = metadata
	})
	ui.createTextForm.AddInputField("Text data:", "", 20, nil, func(textdata string) {
		textRecord.Text = textdata
	})
	ui.createTextForm.AddButton("Create Item", func() {
		if textRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.createTextForm.SetFocus(0)
					ui.pages.ShowPage(NewTextForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectPage)
				})
		}
		if textRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectPage)
			})
		}
		if textRecord.Metadata != "1" && textRecord.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.createTextForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.createTextForm
}

func (ui *UI) readFrm() *tview.Form {
	// loginpsw := ULogin{}
	// ui.readform.AddInputField("Login:", "", 20, nil, func(ulogin string) {
	// 	loginpsw.uLogin = ulogin
	// })
	// ui.readform.AddPasswordField("Password", "", 20, '*', func(upsw string) {
	// 	loginpsw.uPsw = upsw
	// })
	// ui.readform.AddButton("Register", func() {
	// 	if loginpsw.uLogin == "0" {
	// 		log.Print("user reg 0")
	// 		ui.ShowOk("Registration successful. \n Please, keep your credentials safe.", func() {

	// 			ui.pages.SwitchToPage(Menu)
	// 		})
	// 	}
	// 	if loginpsw.uLogin == "1" {
	// 		ui.ShowConfirm("Username already exist", "Do you like try again?",
	// 			func() {
	// 				log.Print("user reg 1")
	// 				ui.readform.SetFocus(0)
	// 				ui.pages.ShowPage(Register)
	// 			},
	// 			func() {
	// 				ui.pages.SwitchToPage(Menu)
	// 			})
	// 	}

	// 	if loginpsw.uLogin != "1" && loginpsw.uLogin != "0" {
	// 		ui.pages.SwitchToPage(Menu)
	// 	}
	// })
	// ui.readform.AddButton("Cancel", func() {
	// 	ui.pages.SwitchToPage(Menu)
	// })
	return ui.readform
}
