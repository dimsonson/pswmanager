package ui

import (
	"strconv"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
)

func (ui *UI) FlexCreate() {
	ui.flexTextCreate = ui.NewCustomFlex(ui.createTextForm, 12)
	ui.flexLoginPairCreate = ui.NewCustomFlex(ui.createLoginPairForm, 12)
	ui.flexBinaryCreate = ui.NewCustomFlex(ui.createBinaryForm, 12)
	ui.flexCardCreate = ui.NewCustomFlex(ui.createCardForm, 15)
	ui.flexSelectCreate = ui.NewCustomFlex(ui.listSelectCreate, 12)
}

func (ui *UI) ListSelectCreate() {
	ui.listSelectCreate = tview.NewList().
		AddItem("Text", "", 'a', func() {
			ui.createTextForm.Clear(true)
			ui.createTextFrm()
			ui.pages.SwitchToPage(NewTextForm)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.createLoginPairForm.Clear(true)
			ui.createLoginPairFrm()
			ui.pages.SwitchToPage(NewLoginPairForm)
		}).
		AddItem("Binary", "", 'c', func() {
			ui.createBinaryForm.Clear(true)
			ui.createBinaryFrm()
			ui.pages.SwitchToPage(NewBinaryForm)
		}).
		AddItem("Card", "", 'd', func() {
			ui.createCardForm.Clear(true)
			ui.createCardFrm()
			ui.pages.SwitchToPage(NewCardForm)
		}).
		AddItem("Main menu", "", 'q', func() {
			ui.pages.SwitchToPage(MainPage)
		})
	ui.listSelectCreate.SetBorder(true)
	ui.listSelectCreate.SetTitle("Create menu")
	ui.listSelectCreate.SetTitleAlign(tview.AlignLeft)
	ui.listSelectCreate.SetWrapAround(true)
	ui.listSelectCreate.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSelectCreate)

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

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if textRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
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

func (ui *UI) createLoginPairFrm() *tview.Form {
	loginPairRecord := models.LoginRecord{}
	ui.createLoginPairForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		loginPairRecord.Metadata = metadata
	})
	ui.createLoginPairForm.AddInputField("Login data:", "", 20, nil, func(logindata string) {
		loginPairRecord.Login = logindata
	})
	ui.createLoginPairForm.AddPasswordField("Password data:", "", 20, '*', func(pswdata string) {
		loginPairRecord.Psw = pswdata
	})
	ui.createLoginPairForm.AddButton("Create Item", func() {
		if loginPairRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new login pair 1")
					ui.createLoginPairForm.SetFocus(0)
					ui.pages.ShowPage(NewLoginPairForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if loginPairRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		if loginPairRecord.Metadata != "1" && loginPairRecord.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.createLoginPairForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.createLoginPairForm
}

func (ui *UI) createBinaryFrm() *tview.Form {
	binaryRecord := models.BinaryRecord{}
	ui.createBinaryForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		binaryRecord.Metadata = metadata
	})
	ui.createBinaryForm.AddInputField("Path to Binary data:", "", 20, nil, func(binarydata string) {
		binaryRecord.Binary = binarydata
	})
	ui.createBinaryForm.AddButton("Create Item", func() {
		if binaryRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.createBinaryForm.SetFocus(0)
					ui.pages.ShowPage(NewTextForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if binaryRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		if binaryRecord.Metadata != "1" && binaryRecord.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.createBinaryForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.createBinaryForm
}

func (ui *UI) createCardFrm() *tview.Form {
	cardRecord := models.CardRecord{}
	var brand = []string{"MIR", "MC", "VISA", "AMEX"}
	ui.createCardForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		cardRecord.Metadata = metadata
	})
	ui.createCardForm.AddDropDown("Brand:", brand, 0, func(branddata string, index int) {
		cardRecord.Brand = models.CardType(index)
	})
	ui.createCardForm.AddInputField("Number:", "", 20, nil, func(numberdata string) {
		cardRecord.Number = numberdata
	})
	ui.createCardForm.AddInputField("Valid Date:", "", 20, nil, func(validdata string) {
		cardRecord.ValidDate = validdata
	})
	ui.createCardForm.AddInputField("Code:", "", 20, nil, func(codedata string) {
		//var err error
		cardRecord.Code, _ = strconv.Atoi(codedata)
	})
	ui.createCardForm.AddInputField("Holder:", "", 20, nil, func(holderdata string) {
		cardRecord.Holder = holderdata
	})
	ui.createCardForm.AddButton("Create Item", func() {
		if cardRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.createCardForm.SetFocus(0)
					ui.pages.ShowPage(NewCardForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if cardRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		if cardRecord.Metadata != "1" && cardRecord.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.createCardForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.createCardForm
}
