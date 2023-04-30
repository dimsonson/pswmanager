package ui

// import (
// 	"os"
// 	"strconv"
// 	"syscall"

// 	"github.com/derailed/tcell/v2"
// 	"github.com/derailed/tview"
// 	"github.com/dimsonson/pswmanager/internal/masterserver/models"
// 	"github.com/dimsonson/pswmanager/pkg/log"
// 	"github.com/rs/zerolog"
// )

// func (ui *UI) FlexRead() {
// 	ui.flexSelectRead = ui.NewAppFlex(ui.listSelectRead, 12)
// 	ui.flexTextRead = ui.NewAppFlex(ui.readTextForm, 10)
// 	ui.flexLoginPairRead = ui.NewAppFlex(ui.readLoginPairForm, 12)
// 	ui.flexBinaryRead = ui.NewAppFlex(ui.readBinaryForm, 12)
// 	ui.flexCardRead = ui.NewAppFlex(ui.readCardForm, 12)
// }

// func (ui *UI) ListSelectRead() {
// 	ui.listSelectRead = tview.NewList().
// 		AddItem("Text Item", "", 'a', func() {
// 			ui.readTextForm.Clear(true)
// 			textRecord:= models.TextRecord{
// 				RecordID: "123",
// 				Metadata: "www.ru",
// 				Text: "secret",
// 			}
// 			ui.readTextFrm(textRecord)
// 			ui.pages.SwitchToPage(ReadTextForm)
// 		}).
// 		AddItem("Login pair", "", 'b', func() {
// 			ui.readLoginPairForm.Clear(true)
// 			ui.readLoginPairFrm()
// 			ui.pages.SwitchToPage(ReadLoginPairForm)
// 		}).
// 		AddItem("Binary", "", 'b', func() {
// 			ui.readBinaryForm.Clear(true)
// 			ui.readBinaryFrm()
// 			ui.pages.SwitchToPage(ReadBinaryForm)
// 		}).
// 		AddItem("Card", "", 'b', func() {
// 			ui.readCardForm.Clear(true)
// 			ui.readCardFrm()
// 			ui.pages.SwitchToPage(ReadCardForm)
// 		}).
// 		AddItem("Quit", "", 'q', func() {
// 			log.Logg = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
// 			ui.MainApp.Stop()
// 			err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
// 			if err != nil {
// 				log.Print("stop programm error")
// 				return
// 			}
// 		})
// 	ui.listSelectRead.SetBorder(true)
// 	ui.listSelectRead.SetTitle("menu")
// 	ui.listSelectRead.SetTitleAlign(tview.AlignLeft)
// 	ui.listSelectRead.SetWrapAround(true)
// 	ui.listSelectRead.SetBackgroundColor(tcell.Color108)
// 	ui.MainApp.SetFocus(ui.listSelectRead)
// }

// func (ui *UI) addTextItems() {
// 	for index, item := range ui.TextSearchResult {
// 		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Text, "", rune(49+index), nil)
// 	}
// }

// func (ui *UI) addLoginItems() {
// 	for index, item := range ui.LoginsSearchResult {
// 		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Login, "", rune(49+index), nil)
// 	}
// }

// func (ui *UI) addBinaryItems() {
// 	for index, item := range ui.BinarySearchResult {
// 		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Binary, "", rune(49+index), nil)
// 	}
// }

// func (ui *UI) addCardItems() {
// 	for index, item := range ui.CardSearchResult {
// 		ui.listTextSearchResult.AddItem(item.Metadata+" "+item.Number, "", rune(49+index), nil)
// 	}
// }

// func (ui *UI) readTextFrm(item models.TextRecord) *tview.Form {
// 	ui.createTextForm.AddInputField("Metadata:", item.Metadata, 20, nil, func(metadata string) {
// 		item.Metadata = metadata
// 	})
// 	ui.createTextForm.AddInputField("Text data:", item.Text, 20, nil, func(textdata string) {
// 		item.Text = textdata
// 	})
// 	ui.createTextForm.AddButton("Create Item", func() {
// 		if item.Metadata == "0" {
// 			ui.ShowConfirm("Error record to database", "Do you like try again?",
// 				func() {
// 					log.Print("new text 1")
// 					ui.readTextForm.SetFocus(0)
// 					ui.pages.ShowPage(ReadTextForm)
// 				},
// 				func() {

// 					ui.pages.SwitchToPage(SelectReadPage)
// 				})
// 		}
// 		if item.Metadata == "1" {
// 			log.Print("user login 0")
// 			ui.ShowOk("New Text Item recorded to database", func() {
// 				ui.pages.SwitchToPage(SelectReadPage)
// 			})
// 		}
// 		if item.Metadata != "1" && item.Metadata != "0" {
// 			ui.pages.SwitchToPage(MainPage)
// 		}
// 	})
// 	ui.createTextForm.AddButton("Cancel", func() {
// 		ui.pages.SwitchToPage(MainPage)
// 	})
// 	return ui.createTextForm
// }

// func (ui *UI) readLoginPairFrm() *tview.Form {
// 	loginPairRecord := models.LoginRecord{}
// 	ui.createLoginPairForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
// 		loginPairRecord.Metadata = metadata
// 	})
// 	ui.createLoginPairForm.AddInputField("Login data:", "", 20, nil, func(logindata string) {
// 		loginPairRecord.Login = logindata
// 	})
// 	ui.createLoginPairForm.AddPasswordField("Password data:", "", 20, '*', func(pswdata string) {
// 		loginPairRecord.Psw = pswdata
// 	})
// 	ui.createLoginPairForm.AddButton("Create Item", func() {
// 		if loginPairRecord.Metadata == "0" {
// 			ui.ShowConfirm("Error record to database", "Do you like try again?",
// 				func() {
// 					log.Print("new login pair 1")
// 					ui.createLoginPairForm.SetFocus(0)
// 					ui.pages.ShowPage(NewLoginPairForm)
// 				},
// 				func() {

// 					ui.pages.SwitchToPage(SelectCreatePage)
// 				})
// 		}
// 		if loginPairRecord.Metadata == "1" {
// 			log.Print("user login 0")
// 			ui.ShowOk("New Text Item recorded to database", func() {
// 				ui.pages.SwitchToPage(SelectCreatePage)
// 			})
// 		}
// 		if loginPairRecord.Metadata != "1" && loginPairRecord.Metadata != "0" {
// 			ui.pages.SwitchToPage(MainPage)
// 		}
// 	})
// 	ui.createLoginPairForm.AddButton("Cancel", func() {
// 		ui.pages.SwitchToPage(MainPage)
// 	})
// 	return ui.createLoginPairForm
// }

// func (ui *UI) readBinaryFrm() *tview.Form {
// 	binaryRecord := models.BinaryRecord{}
// 	ui.createBinaryForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
// 		binaryRecord.Metadata = metadata
// 	})
// 	ui.createBinaryForm.AddInputField("Path to Binary data:", "", 20, nil, func(binarydata string) {
// 		binaryRecord.Binary = binarydata
// 	})
// 	ui.createBinaryForm.AddButton("Create Item", func() {
// 		if binaryRecord.Metadata == "0" {
// 			ui.ShowConfirm("Error record to database", "Do you like try again?",
// 				func() {
// 					log.Print("new text 1")
// 					ui.createBinaryForm.SetFocus(0)
// 					ui.pages.ShowPage(NewTextForm)
// 				},
// 				func() {

// 					ui.pages.SwitchToPage(SelectCreatePage)
// 				})
// 		}
// 		if binaryRecord.Metadata == "1" {
// 			log.Print("user login 0")
// 			ui.ShowOk("New Text Item recorded to database", func() {
// 				ui.pages.SwitchToPage(SelectCreatePage)
// 			})
// 		}
// 		if binaryRecord.Metadata != "1" && binaryRecord.Metadata != "0" {
// 			ui.pages.SwitchToPage(MainPage)
// 		}
// 	})
// 	ui.createBinaryForm.AddButton("Cancel", func() {
// 		ui.pages.SwitchToPage(MainPage)
// 	})
// 	return ui.createBinaryForm
// }

// func (ui *UI) readCardFrm() *tview.Form {
// 	cardRecord := models.CardRecord{}
// 	var brand = []string{"MIR", "MC", "VISA", "AMEX"}
// 	ui.createCardForm.AddInputField("Metadata:", "", 20, nil, func(metadata string) {
// 		cardRecord.Metadata = metadata
// 	})
// 	ui.createCardForm.AddDropDown("Brand:", brand, 0, func(branddata string, index int) {
// 		cardRecord.Brand = models.CardType(index)
// 	})
// 	ui.createCardForm.AddInputField("Number:", "", 20, nil, func(numberdata string) {
// 		cardRecord.Number = numberdata
// 	})
// 	ui.createCardForm.AddInputField("Valid Date:", "", 20, nil, func(validdata string) {
// 		cardRecord.ValidDate = validdata
// 	})
// 	ui.createCardForm.AddInputField("Code:", "", 20, nil, func(codedata string) {
// 		//var err error
// 		cardRecord.Code, _ = strconv.Atoi(codedata)
// 	})
// 	ui.createCardForm.AddInputField("Holder:", "", 20, nil, func(holderdata string) {
// 		cardRecord.Holder = holderdata
// 	})
// 	ui.createCardForm.AddButton("Create Item", func() {
// 		if cardRecord.Metadata == "0" {
// 			ui.ShowConfirm("Error record to database", "Do you like try again?",
// 				func() {
// 					log.Print("new text 1")
// 					ui.createCardForm.SetFocus(0)
// 					ui.pages.ShowPage(NewCardForm)
// 				},
// 				func() {

// 					ui.pages.SwitchToPage(SelectCreatePage)
// 				})
// 		}
// 		if cardRecord.Metadata == "1" {
// 			log.Print("user login 0")
// 			ui.ShowOk("New Text Item recorded to database", func() {
// 				ui.pages.SwitchToPage(SelectCreatePage)
// 			})
// 		}
// 		if cardRecord.Metadata != "1" && cardRecord.Metadata != "0" {
// 			ui.pages.SwitchToPage(MainPage)
// 		}
// 	})
// 	ui.createCardForm.AddButton("Cancel", func() {
// 		ui.pages.SwitchToPage(MainPage)
// 	})
// 	return ui.createCardForm
// }
