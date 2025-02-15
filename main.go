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
	p.SetRect(0, 0, 50, 5)

	sl := widgets.NewSparkline()
	sl.MaxHeight = 2000
	sl.LineColor = ui.ColorGreen
	sl.Data = buffer.GetAll()

	slg := widgets.NewSparklineGroup(sl)
	slg.Title = fmt.Sprintf("Stock Ticker %s", symbol)
	slg.SetRect(0, 5, 50, 15)

	l := widgets.NewList()
	l.Title = "10 last data points"
	l.SetRect(50, 0, 70, 15)

	// Fetch stock data once before entering the loop
	stockInfo, err := stock.GetStockInfo(symbol, apiKey)
	if err != nil {
		p.Text = "Error fetching stock data!"
	} else {
		buffer.Add(float64(*stockInfo.C))

		p.Text = fmt.Sprintf("Stock: %s\nCurrent Price: %f\nChange Percent: %f",
			symbol, *stockInfo.C, *stockInfo.Dp)
		sl.Data = buffer.GetAll()

		lastTen := buffer.GetLastN(10)
		var stringArr []string
		for _, f := range lastTen {
			stringArr = append(stringArr, strconv.FormatFloat(f, 'f', 2, 64))
		}
		l.Rows = stringArr
	}
	ui.Render(p, l, slg)

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

				p.Text = fmt.Sprintf("Stock: %s\nCurrent Price: %f\nChange Percent: %f",
					symbol, *stockInfo.C, *stockInfo.Dp)
				sl.Data = buffer.GetAll()

				lastTen := buffer.GetLastN(10)
				var stringArr []string
				for _, f := range lastTen {
					stringArr = append(stringArr, strconv.FormatFloat(f, 'f', 2, 64))
				}
				l.Rows = stringArr
			}
			ui.Render(p, l, slg)
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
		return 100
	}

	return bufferSize
}
