package stock

import (
	"context"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func GetStockInfo(symbol string, apiKey string) (finnhub.Quote, error) {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	res, _, err := finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()

	return res, err
}
