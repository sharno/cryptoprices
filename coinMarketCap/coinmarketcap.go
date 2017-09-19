package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PossibleConvert lists all the possible currencies to get the prices into
var PossibleConvert = []string{"USD", "AUD", "BRL", "CAD", "CHF", "CLP", "CNY", "CZK", "DKK", "EUR", "GBP",
	"HKD", "HUF", "IDR", "ILS", "INR", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD",
	"PHP", "PKR", "PLN", "RUB", "SEK", "SGD", "THB", "TRY", "TWD", "ZAR"}

// GetPrices gets the prices from CoinMarketCap
func GetPrices(limit uint, convert string) (*[]map[string]string, error) {
	return getAndUnmarshall(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?limit=%v&convert=%s", limit, convert))
}

// GetCoinPrice gets the price of an individual coin
func GetCoinPrice(coin, convert string) (*[]map[string]string, error) {
	return getAndUnmarshall(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/%s/?convert=%s", coin, convert))
}

func getAndUnmarshall(url string) (*[]map[string]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't reach the api of coinmarketcap %s: %v", url, err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the response body sent from %s: %v", url, err)
	}
	data := &[]map[string]string{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, fmt.Errorf("couldn't unmarshall the json sent from %s: %v", url, err)
	}
	return data, nil
}

// IsValidConvert makes sure that the currency converted is a valid one
func IsValidConvert(convert string) bool {
	convertUpperCase := strings.ToUpper(convert)
	for _, v := range PossibleConvert {
		if v == convertUpperCase {
			return true
		}
	}
	return false
}
