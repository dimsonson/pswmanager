package ui

import (
	"encoding/hex"
	"os"
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
					ui.listTextSearchResult.Clear()
					var err error
					ui.TextSearchResult, err = ui.t.SearchText(ui.ctx, ui.searchdata)
					if err != nil {
						log.Print("search text error: ", err)
					}
					if len(ui.TextSearchResult) != 0 {
						ui.addTextItemsList()
						ui.pages.SwitchToPage(TextSearchResult)
					}
					if len(ui.TextSearchResult) == 0 {
						ui.ShowOk("Nothing found for this request.", func() {
						})
					}
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Login pair", "", 'b', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					ui.listLoginsSearchResult.Clear()
					ui.addLoginItemsList()
					var err error
					ui.LoginsSearchResult, err = ui.l.SearchLogin(ui.ctx, ui.searchdata)
					if err != nil {
						log.Print("search logins error: ", err)
					}
					if len(ui.LoginsSearchResult) != 0 {
						ui.addLoginItemsList()
						ui.pages.SwitchToPage(LoginsSearchResult)
					}
					if len(ui.LoginsSearchResult) == 0 {
						ui.ShowOk("Nothing found for this request.", func() {
						})
					}
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Binary", "", 'c', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					ui.listBinarySearchResult.Clear()
					var err error
					ui.BinarySearchResult, err = ui.b.SearchBinary(ui.ctx, ui.searchdata)
					if err != nil {
						log.Print("search binary error: ", err)
					}
					if len(ui.BinarySearchResult) != 0 {
						ui.addBinaryItemsList()
						ui.pages.SwitchToPage(BinarySearchResult)
					}
					if len(ui.BinarySearchResult) == 0 {
						ui.ShowOk("Nothing found for this request.", func() {
						})
					}
				})
			ui.pages.SwitchToPage(searchForm)
		}).
		AddItem("Card", "", 'd', func() {
			ui.searchForm.Clear(true)
			ui.searchFrm(
				func() {
					ui.listCardSearchResult.Clear()
					var err error
					ui.CardSearchResult, err = ui.c.SearchCard(ui.ctx, ui.searchdata)
					if err != nil {
						log.Print("search card error: ", err)
					}
					if len(ui.CardSearchResult) != 0 {
						ui.addCardItemsList()
						ui.pages.SwitchToPage(CardSearchResult)
					}
					if len(ui.CardSearchResult) == 0 {
						ui.ShowOk("Nothing found for this request.", func() {
						})
					}
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
	ui.listTextSearchResult.Clear()
	for index, item := range ui.TextSearchResult {
		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Text, "", rune(49+index), func() {
			ui.readTextForm.Clear(true)
			ui.readTextFrm(item)
			ui.pages.SwitchToPage(ReadTextForm)
		})
	}
}

func (ui *UI) addLoginItemsList() {
	ui.listLoginsSearchResult.Clear()
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
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readTextForm.AddButton("Update Item", func() {
		item.Operation = models.Update
		err := ui.t.ProcessingText(ui.ctx, item)
		if err != nil {
			log.Print("save text data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {
					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readTextForm.AddButton("Delete Item", func() {
		item.Operation = models.Delete
		err := ui.t.ProcessingText(ui.ctx, item)
		if err != nil {
			log.Print("save text data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readTextForm.SetFocus(2)
	ui.readTextForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readTextForm
}

func (ui *UI) readLoginPairFrm(item models.LoginRecord) *tview.Form {
	ui.readLoginPairForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		item.Metadata = metadata
	})
	ui.readLoginPairForm.AddInputField("Login data:", item.Login, 20, nil, func(logindata string) {
		item.Login = logindata
	})
	ui.readLoginPairForm.AddInputField("Password data:", item.Psw, 20, nil, func(pswdata string) {
		item.Psw = pswdata
	})
	ui.readLoginPairForm.AddButton("OK", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readLoginPairForm.AddButton("Update Item", func() {
		item.Operation = models.Update
		err := ui.l.ProcessingLogin(ui.ctx, item)
		if err != nil {
			log.Print("save logins data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {
					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Login Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readLoginPairForm.AddButton("Delete Item", func() {
		item.Operation = models.Delete
		err := ui.l.ProcessingLogin(ui.ctx, item)
		if err != nil {
			log.Print("save login data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readLoginPairForm.SetFocus(3)
	ui.readLoginPairForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readLoginPairForm
}

func (ui *UI) readBinaryFrm(item models.BinaryRecord) *tview.Form {
	ui.readBinaryForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		item.Metadata = metadata
	})
	ui.readBinaryForm.AddInputField("Path to Binary data:", item.Binary, 20, nil, func(binarydata string) {
		item.Binary = binarydata
	})
	ui.readBinaryForm.AddButton("OK", func() {
		binaryFile, err := hex.DecodeString(item.Binary)
		if err != nil {
			log.Print("decode binary file error:", err)
		}
		err = os.WriteFile("binaryfile.tmp", binaryFile, 0666)
		if err != nil {
			log.Print("write binary file error:", err)
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readBinaryForm.AddButton("Update Item", func() {
		item.Operation = models.Update
		err := ui.b.ProcessingBinary(ui.ctx, item)
		if err != nil {
			log.Print("save binary data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {
					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readBinaryForm.AddButton("Delete Item", func() {
		item.Operation = models.Delete
		err := ui.b.ProcessingBinary(ui.ctx, item)
		if err != nil {
			log.Print("save binary data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readBinaryForm.SetFocus(2)
	ui.readBinaryForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.readBinaryForm
}

func (ui *UI) readCardFrm(item models.CardRecord) *tview.Form {
	var brand = []string{"MIR", "MC", "VISA", "AMEX"}
	itemBrand, err := strconv.Atoi(item.Brand)
	if err != nil {
		log.Print("itemBrand conversion error:", err)
	}
	ui.readCardForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
		item.Metadata = metadata
	})
	ui.readCardForm.AddDropDown("Brand:", brand, itemBrand, func(branddata string, index int) {
		item.Brand = strconv.Itoa(index)
	})
	ui.readCardForm.AddInputField("Number:", item.Number, 20, nil, func(numberdata string) {
		item.Number = numberdata
	})
	ui.readCardForm.AddInputField("Valid Date:", item.ValidDate, 20, nil, func(validdata string) {
		item.ValidDate = validdata
	})
	ui.readCardForm.AddInputField("Code:", item.Code, 20, nil, func(codedata string) {
		item.Code = codedata
	})
	ui.readCardForm.AddInputField("Holder:", item.Holder, 20, nil, func(holderdata string) {
		item.Holder = holderdata
	})
	ui.readCardForm.AddButton("OK", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readCardForm.AddButton("Update Item", func() {
		item.Operation = models.Update
		err := ui.c.ProcessingCard(ui.ctx, item)
		if err != nil {
			log.Print("save card data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {
					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
		ui.pages.SwitchToPage(MainPage)
	})
	ui.readCardForm.AddButton("Delete Item", func() {
		item.Operation = models.Delete
		err := ui.c.ProcessingCard(ui.ctx, item)
		if err != nil {
			log.Print("save card data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
				},
				func() {

					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Text Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
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
