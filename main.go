package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	limit := flag.Int("limit", 10, "How many of the top coins on CoinMarketCap to be shown")
	convert := flag.String("Convert", "", "Convert the prices to a different currency other than USD")
	flag.Parse()
	PrintPrices(*limit, *convert)
}

// Currency is the type of response back from CoinMarketCap
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

// GetPrices gets the prices from CoinMarketCap
func GetPrices(limit int, convert string) (*[]Currency, error) {
	res, err := http.Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?limit=%v&convert=%s", limit, convert))
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

// PrintPrices gets the prices from CoinMarketCap and prints them in the console
func PrintPrices(limit int, convert string) {
	currencies, err := GetPrices(limit, convert)
	if err != nil {
		log.Printf("couldn't get the prices of cryptocurrencies: %v", err)
		return
	}
	for _, currency := range *currencies {
		fmt.Printf("%s: %v USD\n", currency.Name, currency.PriceUsd)
	}
}
