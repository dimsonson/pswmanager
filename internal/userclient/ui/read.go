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

	ui.flexTextSearchResult = ui.NewCustomFlex(ui.listTextSearchResult, 15)
	ui.flexBinarySearchResult = ui.NewCustomFlex(ui.listBinarySearchResult, 15)
	ui.flexCardSearchResult = ui.NewCustomFlex(ui.listCardSearchResult, 15)
	ui.flexLoginPairSearchResult = ui.NewCustomFlex(ui.listLoginsSearchResult, 15)

	ui.flexSearchForm = ui.NewCustomFlex(ui.searchForm, 10)
}

func (ui *UI) ListSelectRead() {
	ui.listSelectRead = tview.NewList().
		AddItem("Text Item", "", 'a', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					log.Print("searchdata: ", ui.searchdata)
					ui.listTextSearchResult.Clear()
					ui.addTextItemsList()
					log.Print("addTextItems: ", ui.listTextSearchResult)
					ui.pages.SwitchToPage(TextSearchResult)
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					log.Print("searchdata ", ui.searchdata)
					ui.listLoginsSearchResult.Clear()
					ui.addLoginItemsList()
					ui.pages.SwitchToPage(LoginsSearchResult)
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Binary", "", 'c', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					log.Print("searchdata", ui.searchdata)
					ui.listBinarySearchResult.Clear()
					ui.addBinaryItemsList()
					ui.pages.SwitchToPage(BinarySearchResult)
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Card", "", 'd', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					log.Print("searchdata", ui.searchdata)
					ui.listCardSearchResult.Clear()
					ui.addCardItemsList()
					ui.pages.SwitchToPage(CardSearchResult)
				})
			ui.pages.SwitchToPage(searchForm)
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

func (ui *UI) addTextItemsList() {
	//ui.listTextSearchResult.Clear()
	for index, item := range ui.TextSearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Text, "", rune(49+index), func() {
			ui.readTextForm.Clear(true)
			ui.readTextFrm(item)
			ui.pages.SwitchToPage(ReadTextForm)
		})
	}
}

func (ui *UI) addLoginItemsList() {
	//ui.listLoginsSearchResult.Clear()
	for index, item := range ui.LoginsSearchResult {
		ui.listLoginsSearchResult.AddItem(item.Metadata+" "+item.Login, "", rune(49+index), func() {
			ui.readLoginPairForm.Clear(true)
			ui.readLoginPairFrm(item)
			ui.pages.SwitchToPage(ReadLoginPairForm)
		})
	}
}

func (ui *UI) addBinaryItemsList() {
	ui.listBinarySearchResult.Clear()
	for index, item := range ui.BinarySearchResult {
		ui.listBinarySearchResult.AddItem(item.Metadata+" "+item.Metadata, "", rune(49+index), func() {
			ui.readBinaryForm.Clear(true)
			ui.readBinaryFrm(item)
			ui.pages.SwitchToPage(ReadBinaryForm)
		})
	}
}

func (ui *UI) addCardItemsList() {
	ui.listCardSearchResult.Clear()
	for index, item := range ui.CardSearchResult {
		ui.listCardSearchResult.AddItem(item.Metadata+" "+item.Number, "", rune(49+index), func() {
			ui.readCardForm.Clear(true)
			ui.readCardFrm(item)
			ui.pages.SwitchToPage(ReadCardForm)
		})
	}
}

func (ui *UI) readTextFrm(item models.TextRecord) *tview.Form {
	ui.readTextForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		item.Metadata = metadata
	})
	ui.readTextForm.AddInputField("Text data:", item.Text, 20, nil, func(textdata string) {
		item.Text = textdata
	})
	ui.readTextForm.AddButton("OK", func() {
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
	ui.readTextForm.AddButton("Update Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readTextForm.AddButton("Delete Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readTextForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readTextForm.SetFocus(2)
	ui.readTextForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readTextForm
}

func (ui *UI) readLoginPairFrm(item models.LoginRecord) *tview.Form {
	loginPairRecord := models.LoginRecord{}
	ui.readLoginPairForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		loginPairRecord.Metadata = metadata
	})
	ui.readLoginPairForm.AddInputField("Login data:", item.Login, 20, nil, func(logindata string) {
		loginPairRecord.Login = logindata
	})
	ui.readLoginPairForm.AddPasswordField("Password data:", item.Psw, 20, '*', func(pswdata string) {
		loginPairRecord.Psw = pswdata
	})
	ui.readLoginPairForm.AddButton("OK", func() {
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
	ui.readLoginPairForm.AddButton("Update Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readLoginPairForm.AddButton("Delete Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readLoginPairForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readLoginPairForm.SetFocus(3)

	ui.readLoginPairForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readLoginPairForm
}

func (ui *UI) readBinaryFrm(item models.BinaryRecord) *tview.Form {
	binaryRecord := models.BinaryRecord{}
	ui.readBinaryForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		binaryRecord.Metadata = metadata
	})
	ui.readBinaryForm.AddInputField("Path to Binary data:", item.Binary, 20, nil, func(binarydata string) {
		binaryRecord.Binary = binarydata
	})
	ui.readBinaryForm.AddButton("OK", func() {
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

	ui.readBinaryForm.AddButton("Update Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readBinaryForm.AddButton("Delete Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readBinaryForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readBinaryForm.SetFocus(2)

	ui.readBinaryForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readBinaryForm
}

func (ui *UI) readCardFrm(item models.CardRecord) *tview.Form {
	cardRecord := models.CardRecord{}
	var brand = []string{"MIR", "MC", "VISA", "AMEX"}
	ui.readCardForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		cardRecord.Metadata = metadata
	})
	ui.readCardForm.AddDropDown("Brand:", brand, int(item.Brand), func(branddata string, index int) {
		cardRecord.Brand = models.CardType(index)
	})
	ui.readCardForm.AddInputField("Number:", item.Number, 20, nil, func(numberdata string) {
		cardRecord.Number = numberdata
	})
	ui.readCardForm.AddInputField("Valid Date:", item.ValidDate, 20, nil, func(validdata string) {
		cardRecord.ValidDate = validdata
	})
	ui.readCardForm.AddInputField("Code:", "001_int", 20, nil, func(codedata string) {
		//var err error
		cardRecord.Code, _ = strconv.Atoi(codedata)
	})
	ui.readCardForm.AddInputField("Holder:", item.Holder, 20, nil, func(holderdata string) {
		cardRecord.Holder = holderdata
	})
	ui.readCardForm.AddButton("OK", func() {
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

	ui.readCardForm.AddButton("Update Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readCardForm.AddButton("Delete Item", func() {
		ui.pages.SwitchToPage(MainPage)
	})

	ui.readCardForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readCardForm.SetFocus(6)
	return ui.readCardForm
}

func (ui *UI) searchFrm(ack confirmFunc) *tview.Form {
	ui.searchForm.AddInputField("Enter Metadata for search:", "", 20, nil, func(sdata string) {
		ui.searchdata = sdata
	})
	ui.searchForm.AddButton("Search now", ack)
	ui.searchForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.searchForm
}
