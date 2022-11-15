package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gold/repository"
	"strconv"
	"time"
)

func (app *Config) getToolBar(win fyne.Window) *widget.Toolbar {
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addHoldingsDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.refreshPriceContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			win := app.showReferences()
			win.Resize(fyne.NewSize(320, 20))
			win.Show()
		}),
	)
	return toolBar
}
func (app *Config) addHoldingsDialog() dialog.Dialog {
	//new entry
	addAmountEntry := widget.NewEntry()
	purchaseDateEntry := widget.NewEntry()
	purchasePriceEntry := widget.NewEntry()
	//entry set validator
	addAmountEntry.Validator = validation.NewRegexp("^[1-9]\\d*$", "please input valid number")
	purchasePriceEntry.Validator = validation.NewRegexp("^\\d+(\\.\\d{1,2})?$", "please input a valid price")
	timeRegex := "(([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|" +
		"((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|" +
		"((0[48]|[2468][048]|[3579][26])00))-02-29)$"
	purchaseDateEntry.Validator = validation.NewRegexp(timeRegex, "please write a correct date")
	//set ref of entry
	app.AddHoldingsPurchaseAmountEntry = addAmountEntry
	app.AddHoldingsPurchaseDateEntry = purchaseDateEntry
	app.AddHoldingsPurchasePriceEntry = purchasePriceEntry
	purchaseDateEntry.PlaceHolder = "YYYY-MM-DD"
	//create a dialog
	addForm := dialog.NewForm(
		"Add Gold",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Amount int toz", Widget: addAmountEntry},
			{Text: "Purchase Price", Widget: purchasePriceEntry},
			{Text: "Purchase Date", Widget: purchaseDateEntry},
		},
		func(valid bool) {
			if !valid {
				return
			}
			amount, _ := strconv.Atoi(addAmountEntry.Text)
			purchasePrice, _ := strconv.ParseFloat(purchasePriceEntry.Text, 32)
			purchaseDate, _ := time.Parse("2006-01-02", purchaseDateEntry.Text)
			purchasePrice *= 100
			_, err := app.DB.InsertHolding(repository.Holdings{
				Amount:        amount,
				PurchaseDate:  purchaseDate,
				PurchasePrice: int(purchasePrice),
			})
			if err != nil {
				app.ErrorLog.Println(err)
			}
			app.refreshHoldingsTable()
		},
		app.MainWindow,
	)
	//size and show the dialog
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()
	return addForm
}
func (app *Config) showReferences() fyne.Window {
	//create a new window
	win := app.App.NewWindow("Preferences")
	//create a label
	lbl := widget.NewLabel("Preferred currency")
	//create a currency select container
	cur := widget.NewSelect([]string{"USD", "CAD", "GBP"}, func(s string) {
		currency = s
		app.App.Preferences().SetString("currency", s)
	})
	cur.Selected = currency
	//create a save button
	btn := widget.NewButton("Save", func() {
		win.Close()
		app.refreshPriceContent()
	})
	btn.Importance = widget.HighImportance
	//set content for window
	win.SetContent(container.NewVBox(lbl, cur, btn))
	return win
}
