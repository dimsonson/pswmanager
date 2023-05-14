package ui

import (
	"context"
	"strconv"
	"time"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/google/uuid"
)

type TextServicesProvider interface {
	ProcessingText(ctx context.Context, record models.TextRecord) error
	SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error)
}

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
		textRecord.UID = ui.cfg.UserID
		textRecord.AppID = ui.cfg.AppID
		textRecord.Operation = models.Create
		textRecord.ChngTime = time.Now()
		textRecord.RecordID = uuid.NewString()
		err := ui.t.ProcessingText(ui.ctx, textRecord)
		if err != nil {
			log.Print("save text data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
					ui.pages.ShowPage(NewTextForm)
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
		loginPairRecord.UID = ui.cfg.UserID
		loginPairRecord.AppID = ui.cfg.AppID
		loginPairRecord.Operation = models.Create
		loginPairRecord.ChngTime = time.Now()
		loginPairRecord.RecordID = uuid.NewString()
		err := ui.l.ProcessingLogin(ui.ctx, loginPairRecord)
		if err != nil {
			log.Print("save login data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
					ui.pages.ShowPage(NewLoginPairForm)
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
		binaryRecord.UID = ui.cfg.UserID
		binaryRecord.AppID = ui.cfg.AppID
		binaryRecord.Operation = models.Create
		binaryRecord.ChngTime = time.Now()
		binaryRecord.RecordID = uuid.NewString()
		err := ui.b.ProcessingBinary(ui.ctx, binaryRecord)
		if err != nil {
			log.Print("save login data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
					ui.pages.ShowPage(NewBinaryForm)
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
		cardRecord.UID = ui.cfg.UserID
		cardRecord.AppID = ui.cfg.AppID
		cardRecord.Operation = models.Create
		cardRecord.ChngTime = time.Now()
		cardRecord.RecordID = uuid.NewString()
		err := ui.c.ProcessingCard(ui.ctx, cardRecord)
		if err != nil {
			log.Print("save card data error:", err)
			ui.ShowConfirm("Error record to database", "Do you like try again?",
				func() {
					ui.createTextForm.SetFocus(0)
					ui.pages.ShowPage(NewLoginPairForm)
				},
				func() {
					ui.pages.SwitchToPage(SelectCreatePage)
				})
		}
		if err == nil {
			ui.ShowOk("New Card Item recorded to database", func() {
				ui.pages.SwitchToPage(SelectCreatePage)
			})
		}
	})
	ui.createCardForm.AddButton("Cancel", func() {
		ui.pages.SwitchToPage(MainPage)
	})
	return ui.createCardForm
}
