package main

import (
	"bytes"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/container"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func (app *Config) pricesTab() *fyne.Container {
	chart := app.getChart()
	box := container.NewVBox(chart)
	box.Resize(fyne.NewSize(770, 410))
	app.PriceChartContainer = box
	return box
}
func (app *Config) getChart() *canvas.Image {
	apiURL := fmt.Sprintf("https://goldprice.org/charts/gold_3d_b_o_%s_x.png", strings.ToLower(currency))
	var img *canvas.Image
	err := app.downloadFile(apiURL, "gold.png")
	if err != nil {
		// use bundled image
		img = canvas.NewImageFromResource(resourceUnreachablePng)
	} else {
		img = canvas.NewImageFromFile("gold.png")
	}
	img.SetMinSize(fyne.NewSize(770, 410))
	img.FillMode = canvas.ImageFillOriginal
	return img
}
func (app Config) downloadFile(URL, fileName string) error {
	//get the response bytes form calling an url
	resp, err := app.HTTPClient.Get(URL)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("received wrong response code when downloading image")
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	out, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
func (app *Config) priceRefreshTicker() {
	for {
		select {
		case <-time.Tick(30 * time.Second):
			app.refreshPriceContent()
		}
	}
}
