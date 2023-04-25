package ui

import (
	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
)

func (ui *UI) ShowConfirm(title, msg string, ack confirmFunc, cancel cancelFunc) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.Color110).
		SetButtonTextColor(tcell.Color111).
		SetLabelColor(tcell.Color112.TrueColor()).
		SetFieldTextColor(tcell.Color114.TrueColor())
	f.AddButton("Cancel", func() {
		ui.dismiss(ui.Pages)
		cancel()
	})
	f.AddButton("OK", func() {
		ui.dismiss(ui.Pages)
		ack()
	})
	for i := 0; i < 2; i++ {
		b := f.GetButton(i)
		if b == nil {
			continue
		}
		b.SetBackgroundColorActivated(tcell.Color116.TrueColor())
		b.SetLabelColorActivated(tcell.Color117.TrueColor())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm(title, f)
	modal.SetText(msg)
	modal.SetTextColor(tcell.Color119.TrueColor())
	modal.SetDoneFunc(func(int, string) {
		ui.dismiss(ui.Pages)
		cancel()
	})
	ui.pages.AddPage("dialogKey", modal, false, false)
	ui.pages.ShowPage("dialogKey")
}

func (ui *UI) dismiss(pages Pages) {
	ui.pages.RemovePage("dialogKey")
}

func (ui *UI) ShowOk(msg string, ok confirmFunc) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.Color110).
		SetButtonTextColor(tcell.Color111).
		SetLabelColor(tcell.Color112.TrueColor()).
		SetFieldTextColor(tcell.Color114.TrueColor())
	f.AddButton("OK", func() {
		ui.dismiss(ui.Pages)
		ok()
	})
	if b := f.GetButton(0); b != nil {
		b.SetBackgroundColorActivated(tcell.Color116.TrueColor())
		b.SetLabelColorActivated(tcell.Color117.TrueColor())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("", f)
	modal.SetText(msg)
	modal.SetTextColor(tcell.ColorOrangeRed)
	modal.SetDoneFunc(func(int, string) {
		ui.dismiss(ui.Pages)
		ok()
	})
	ui.pages.AddPage("dialogKey", modal, false, false)
	ui.pages.ShowPage("dialogKey")
}
