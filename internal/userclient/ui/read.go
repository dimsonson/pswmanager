package ui

import (
	"strconv"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
)

func (ui *UI) FlexRead() {
	ui.flexSelectRead = ui.NewCustomFlex(ui.listSelectRead, 12)
	ui.flexTextRead = ui.NewCustomFlex(ui.readTextForm, 10)
	ui.flexLoginPairRead = ui.NewCustomFlex(ui.readLoginPairForm, 12)
	ui.flexBinaryRead = ui.NewCustomFlex(ui.readBinaryForm, 12)
	ui.flexCardRead = ui.NewCustomFlex(ui.readCardForm, 15)

	ui.flexTextReadSearch = ui.NewCustomFlex(ui.listTextSearchResult, 15)
}

func (ui *UI) ListSelectRead() {
	ui.listSelectRead = tview.NewList().
		AddItem("Text Item", "", 'a', func() {
			ui.listTextSearchResult.Clear()
			// textRecord := models.TextRecord{
			// 	RecordID: "123",
			// 	Metadata: "www.ru",
			// 	Text:     "secret",
			// }
			//ui.readTextFrm(textRecord)
			ui.addTextItems()
			ui.pages.SwitchToPage(ReadBinarySearch)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.readLoginPairForm.Clear(true)
			ui.readLoginPairFrm()
			ui.pages.SwitchToPage(ReadLoginPairForm)
		}).
		AddItem("Binary", "", 'c', func() {
			ui.readBinaryForm.Clear(true)
			ui.readBinaryFrm()
			ui.pages.SwitchToPage(ReadBinaryForm)
		}).
		AddItem("Card", "", 'd', func() {
			ui.readCardForm.Clear(true)
			ui.readCardFrm()
			ui.pages.SwitchToPage(ReadCardForm)
		}).
		AddItem("Main menu", "", 'q', func() {
			ui.pages.SwitchToPage(MainPage)
		})
	ui.listSelectRead.SetBorder(true)
	ui.listSelectRead.SetTitle("Read menu")
	ui.listSelectRead.SetTitleAlign(tview.AlignLeft)
	ui.listSelectRead.SetWrapAround(true)
	ui.listSelectRead.SetBackgroundColor(tcell.Color108)
	ui.MainApp.SetFocus(ui.listSelectRead)
}

func (ui *UI) addTextItems() {
	for index, item := range ui.TextSearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Text, "", rune(49+index), func() {
			ui.readTextForm.Clear(true)
			ui.readTextFrm(item)
			ui.pages.SwitchToPage(ReadTextForm)
		})
	}
}

func (ui *UI) addLoginItems() {
	for index, item := range ui.LoginsSearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Login, "", rune(49+index), nil)
	}
}

func (ui *UI) addBinaryItems() {
	for index, item := range ui.BinarySearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Binary, "", rune(49+index), nil)
	}
}

func (ui *UI) addCardItems() {
	for index, item := range ui.CardSearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Number, "", rune(49+index), nil)
	}
}

func (ui *UI) readTextFrm(item models.TextRecord) *tview.Form {
	ui.readTextForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		item.Metadata = metadata
	})
	ui.readTextForm.AddInputField("Text data:", item.Text, 20, nil, func(textdata string) {
		item.Text = textdata
	})
	ui.readTextForm.AddButton("Create Item", func() {
		if item.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.readTextForm.SetFocus(0)
					ui.pages.ShowPage(ReadTextForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectReadPage)
				})
		}
		if item.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectReadPage)
			})
		}
		if item.Metadata != "1" && item.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.readTextForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readTextForm
}

func (ui *UI) readLoginPairFrm() *tview.Form {
	loginPairRecord := models.LoginRecord{}
	ui.readLoginPairForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		loginPairRecord.Metadata = metadata
	})
	ui.readLoginPairForm.AddInputField("Login data:", "", 20, nil, func(logindata string) {
		loginPairRecord.Login = logindata
	})
	ui.readLoginPairForm.AddPasswordField("Password data:", "", 20, '*', func(pswdata string) {
		loginPairRecord.Psw = pswdata
	})
	ui.readLoginPairForm.AddButton("Create Item", func() {
		if loginPairRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new login pair 1")
					ui.readLoginPairForm.SetFocus(0)
					ui.pages.ShowPage(ReadLoginPairForm)
				},
				func() {

					ui.pages.SwitchToPage(SelectReadPage)
				})
		}
		if loginPairRecord.Metadata == "1" {
			log.Print("user login 0")
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectReadPage)
			})
		}
		if loginPairRecord.Metadata != "1" && loginPairRecord.Metadata != "0" {
			ui.pages.SwitchToPage(MainPage)
		}
	})
	ui.readLoginPairForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readLoginPairForm
}

func (ui *UI) readBinaryFrm() *tview.Form {
	binaryRecord := models.BinaryRecord{}
	ui.readBinaryForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		binaryRecord.Metadata = metadata
	})
	ui.readBinaryForm.AddInputField("Path to Binary data:", "", 20, nil, func(binarydata string) {
		binaryRecord.Binary = binarydata
	})
	ui.readBinaryForm.AddButton("Create Item", func() {
		if binaryRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.readBinaryForm.SetFocus(0)
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
	ui.readBinaryForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readBinaryForm
}

func (ui *UI) readCardFrm() *tview.Form {
	cardRecord := models.CardRecord{}
	var brand = []string{"MIR", "MC", "VISA", "AMEX"}
	ui.readCardForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
		cardRecord.Metadata = metadata
	})
	ui.readCardForm.AddDropDown("Brand:", brand, 0, func(branddata string, index int) {
		cardRecord.Brand = models.CardType(index)
	})
	ui.readCardForm.AddInputField("Number:", "", 20, nil, func(numberdata string) {
		cardRecord.Number = numberdata
	})
	ui.readCardForm.AddInputField("Valid Date:", "", 20, nil, func(validdata string) {
		cardRecord.ValidDate = validdata
	})
	ui.readCardForm.AddInputField("Code:", "", 20, nil, func(codedata string) {
		//var err error
		cardRecord.Code, _ = strconv.Atoi(codedata)
	})
	ui.readCardForm.AddInputField("Holder:", "", 20, nil, func(holderdata string) {
		cardRecord.Holder = holderdata
	})
	ui.readCardForm.AddButton("Create Item", func() {
		if cardRecord.Metadata == "0" {
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					log.Print("new text 1")
					ui.readCardForm.SetFocus(0)
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
	ui.readCardForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readCardForm
}
