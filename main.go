package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"gold/repository"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App                            fyne.App
	InfoLog                        *log.Logger
	ErrorLog                       *log.Logger
	DB                             *repository.SQLiteRepository
	MainWindow                     fyne.Window
	PriceContainer                 *fyne.Container
	PriceChartContainer            *fyne.Container
	Holdings                       [][]interface{}
	HoldingTable                   *widget.Table
	HTTPClient                     *http.Client
	AddHoldingsPurchaseAmountEntry *widget.Entry
	AddHoldingsPurchaseDateEntry   *widget.Entry
	AddHoldingsPurchasePriceEntry  *widget.Entry
}

var myApp Config

func main() {
	//create a fyne application
	fyneApp := app.NewWithID("com.gocode.goldwatcher.www")
	myApp.App = fyneApp
	myApp.HTTPClient = &http.Client{}
	//create our logger
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//open a connection to the database
	path := "./repository/data/sql.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		myApp.ErrorLog.Println(err)
	}
	myApp.DB = repository.NewSQLiteRepository(db)
	//	create a database repository
	err = myApp.DB.Migrate()
	if err != nil {
		fmt.Println("create a database repository failed:", err)
	}
	//create and size a fyne window
	myApp.MainWindow = fyneApp.NewWindow("GoldWatcher")
	myApp.MainWindow.Resize(fyne.NewSize(770, 410))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()
	myApp.makeUI()
	//show and run the application
	myApp.MainWindow.ShowAndRun()
}
