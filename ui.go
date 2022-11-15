package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) makeUI() {
	//get the current price of gold
	openPrice, currentPrice, priceChange := app.getPriceText()
	//put some information into a container
	priceContent := container.NewGridWithColumns(3, openPrice, currentPrice, priceChange)
	app.PriceContainer = priceContent
	//set interval to refresh priceContainer
	go app.priceRefreshTicker()
	//get toolbar
	toolBar := app.getToolBar(app.MainWindow)
	//get app tabs
	pricesTabContent := app.pricesTab()
	holdingsTabContent := app.holdingsTab()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Prices", theme.HomeIcon(), pricesTabContent),
		container.NewTabItemWithIcon("Holdings", theme.InfoIcon(), holdingsTabContent),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	//add container to window
	finalContent := container.NewVBox(priceContent, toolBar, tabs)
	app.MainWindow.SetContent(finalContent)
}
func (app *Config) refreshPriceContent() {
	app.InfoLog.Println("refresh priceTab")
	open, current, change := app.getPriceText()
	app.PriceContainer.Objects = []fyne.CanvasObject{
		open, current, change,
	}
	app.PriceContainer.Refresh()
	chart := app.getChart()
	app.PriceChartContainer.Objects = []fyne.CanvasObject{
		chart,
	}
	app.PriceChartContainer.Refresh()
}
func (app *Config) refreshHoldingsTable() {
	app.Holdings = app.getHoldingsSlice()
	app.HoldingTable.Refresh()
}
