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

	"termstockticker/stock"
)

func main() {
	apiKey := getApiKey()
	symbol := getStockSymbol()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.SetRect(0, 0, 50, 5)

	for {
		stockInfo, err := stock.GetStockInfo(symbol, apiKey)
		if err != nil {
			p.Text = "Error fetching stock data!"
		} else {
			p.Text = fmt.Sprintf("Stock: %s\nPrice: %f\nChange Percent: %f",
				symbol, *stockInfo.C, *stockInfo.Dp)
		}
		ui.Render(p)

		select {
		case e := <-ui.PollEvents():
			if e.Type == ui.KeyboardEvent && e.ID == "q" {
				return
			}
		case <-time.After(15 * time.Second):
		}
	}
}

func getApiKey() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("API_KEY")
}

func getStockSymbol() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("STOCK_SYMBOL")
}

func getBufferSize() int {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	bufferSizeStr := os.Getenv("BUFFER_SIZE")
	bufferSize, err := strconv.Atoi(bufferSizeStr)

	if err != nil {
		return 100
	}

	return bufferSize
}
