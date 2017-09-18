package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	possibleConvert := []string{"AUD", "BRL", "CAD", "CHF", "CLP", "CNY", "CZK", "DKK", "EUR", "GBP",
		"HKD", "HUF", "IDR", "ILS", "INR", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD",
		"PHP", "PKR", "PLN", "RUB", "SEK", "SGD", "THB", "TRY", "TWD", "ZAR"}

	limit := flag.Int("limit", 10, "How many of the top coins on CoinMarketCap to be shown")
	convert := flag.String("convert", "USD",
		fmt.Sprintf("Convert the prices to a different currency other than USD, possible currencies are: %v", possibleConvert))
	flag.Parse()

	convertLowCase := strings.ToLower(*convert)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go PrintPricesForever(*limit, convertLowCase, wg)
	wg.Wait()
}

// GetPrices gets the prices from CoinMarketCap
func GetPrices(limit int, convert string) (*[]map[string]string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?limit=%v&convert=%s", limit, convert))
	if err != nil {
		return nil, fmt.Errorf("couldn't reach the api for prices: %v", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the response body sent from server: %v", err)
	}
	currencies := &[]map[string]string{}
	if err := json.Unmarshal(b, currencies); err != nil {
		return nil, fmt.Errorf("couldn't unmarshall the json sent from the server: %v", err)
	}
	return currencies, nil
}

// PrintPricesForever gets the prices from CoinMarketCap and prints them in the console
func PrintPricesForever(limit int, convert string, wg *sync.WaitGroup) {
	fmt.Println("CoinMarketCap changes the price every 5 minutes. The new price would be printed here every 5 minutes")
	currencies, err := GetPrices(limit, convert)
	if err != nil {
		log.Printf("couldn't get the prices of cryptocurrencies: %v", err)
		return
	}
	for _, currency := range *currencies {
		fmt.Printf("%s: %v %s\n", currency["name"], currency[fmt.Sprintf("price_%s", convert)], strings.ToUpper(convert))
	}

	ticker := time.Tick(5 * time.Minute)
	for _ = range ticker {
		currencies, err := GetPrices(limit, convert)
		if err != nil {
			log.Printf("couldn't get the prices of cryptocurrencies: %v", err)
			return
		}
		for _, currency := range *currencies {
			fmt.Printf("%s: %v %s\n", currency["name"], currency[fmt.Sprintf("price_%s", convert)], strings.ToUpper(convert))
		}
	}

	wg.Done()
}
