package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	PrintPrices()
}

type Currency struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUsd         float64 `json:"price_usd,string"`
	PriceBtc         float64 `json:"price_btc,string"`
	VolumeUsd24h     float64 `json:"24h_volume_usd,string"`
	MarketCapUsd     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1H  float64 `json:"percent_change_1h,string"`
	PercentChange24H float64 `json:"percent_change_24h,string"`
	PercentChange7D  float64 `json:"percent_change_7d,string"`
	LastUpdated      int64   `json:"last_updated,string"`
	PriceEur         float64 `json:"price_eur,string"`
	VolumeEur24h     float64 `json:"24h_volume_eur,string"`
	MarketCapEur     float64 `json:"market_cap_eur,string"`
}

func GetPrices() (*[]Currency, error) {
	res, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		return nil, fmt.Errorf("couldn't reach the api for prices: %v", err)
	}
	defer res.Body.Close()
	currencies := &[]Currency{}
	if err = json.NewDecoder(res.Body).Decode(currencies); err != nil {
		return nil, fmt.Errorf("couldn't decode json of prices: %v", err)
	}

	return currencies, nil
}

func PrintPrices() {
	currencies, err := GetPrices()
	if err != nil {
		log.Printf("couldn't get the prices of cryptocurrencies: %v", err)
		return
	}
	for _, currency := range *currencies {
		fmt.Printf("%s: %v USD\n", currency.Name, currency.PriceUsd)
	}
}
