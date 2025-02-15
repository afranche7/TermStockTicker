package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/joho/godotenv"

	"termstockticker/circularbuffer"
	"termstockticker/stock"
)

func main() {
	loadEnvironmentVariables()
	apiKey := getApiKey()
	symbol := getStockSymbol()
	bufferSize := getBufferSize()

	buffer := circularbuffer.CreateBuffer(bufferSize)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = fmt.Sprintf("Stock Ticker %s", symbol)
	p.BorderStyle.Fg = ui.ColorCyan
	p.SetRect(0, 0, 60, 4)

	lc := widgets.NewPlot()
	lc.Marker = widgets.MarkerDot
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorGreen
	lc.DrawDirection = widgets.DrawLeft
	lc.BorderStyle.Fg = ui.ColorCyan
	lc.SetRect(0, 4, 60, 15)

	l := widgets.NewList()
	l.Title = "Latest data"
	l.BorderStyle.Fg = ui.ColorCyan
	l.SetRect(60, 0, 75, 15)

	// Fetch stock data once before entering the loop
	stockInfo, err := stock.GetStockInfo(symbol, apiKey)
	if err != nil {
		p.Text = "Error fetching stock data!"
	} else {
		buffer.Add(float64(*stockInfo.C))

		p.Text = fmt.Sprintf("Current Price: %f\nChange Percent: %f", *stockInfo.C, *stockInfo.Dp)
		lc.Data = [][]float64{buffer.GetAll()}

		lastTen := buffer.GetLastN(15)
		var stringArr []string
		for _, f := range lastTen {
			stringArr = append(stringArr, strconv.FormatFloat(f, 'f', 2, 64))
		}
		l.Rows = stringArr
	}
	ui.Render(p, l, lc)

	for {
		select {
		case e := <-ui.PollEvents():
			if e.Type == ui.KeyboardEvent && e.ID == "q" {
				return
			}
		case <-time.After(15 * time.Second):
			stockInfo, err := stock.GetStockInfo(symbol, apiKey)
			if err != nil {
				p.Text = "Error fetching stock data!"
			} else {
				buffer.Add(float64(*stockInfo.C))

				p.Text = fmt.Sprintf("Current Price: %f\nChange Percent: %f", *stockInfo.C, *stockInfo.Dp)
				lc.Data = [][]float64{buffer.GetAll()}

				lastTen := buffer.GetLastN(10)
				var stringArr []string
				for _, f := range lastTen {
					stringArr = append(stringArr, strconv.FormatFloat(f, 'f', 2, 64))
				}
				l.Rows = stringArr
			}
			ui.Render(p, l, lc)
		}
	}
}

func loadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getApiKey() string {
	return os.Getenv("API_KEY")
}

func getStockSymbol() string {
	return os.Getenv("STOCK_SYMBOL")
}

func getBufferSize() int {
	bufferSizeStr := os.Getenv("BUFFER_SIZE")
	bufferSize, err := strconv.Atoi(bufferSizeStr)

	if err != nil {
		return 50
	}

	return bufferSize
}
