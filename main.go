package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	possibleConvert := []string{"AUD", "BRL", "CAD", "CHF", "CLP", "CNY", "CZK", "DKK", "EUR", "GBP",
		"HKD", "HUF", "IDR", "ILS", "INR", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD",
		"PHP", "PKR", "PLN", "RUB", "SEK", "SGD", "THB", "TRY", "TWD", "ZAR"}

	limit := flag.Int("limit", 2, "How many of the top coins on CoinMarketCap to be shown")
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
	coins := &[]map[string]string{}
	if err := json.Unmarshal(b, coins); err != nil {
		return nil, fmt.Errorf("couldn't unmarshall the json sent from the server: %v", err)
	}
	return coins, nil
}

// PrintPricesForever gets the prices from CoinMarketCap and prints them in the console every 5 minutes
func PrintPricesForever(limit int, convert string, wg *sync.WaitGroup) {
	fmt.Printf("CoinMarketCap changes the price every 5 minutes. The new price in %s would be printed here every 5 minutes\n\n",
		strings.ToUpper(convert))

	prevPrices := &map[string]float64{}
	PrintPrices(limit, convert, time.Now(), true, prevPrices)

	ticker := time.Tick(5 * time.Minute)
	for t := range ticker {
		PrintPrices(limit, convert, t, false, prevPrices)
	}

	wg.Done()
}

// PrintPrices gets the prices from CoinMarketCap and prints them in the console every 5 minutes
func PrintPrices(limit int, convert string, t time.Time, withNames bool, prevPrices *map[string]float64) {
	coins, err := GetPrices(limit, convert)
	if err != nil {
		log.Printf("couldn't get the prices of coins: %v", err)
		return
	}
	if withNames {
		fmt.Printf("%-18s", "timestamp")
		for _, coin := range *coins {
			fmt.Printf("%-14s%-14s| ", coin["name"], "diff")
		}
		fmt.Println()
	}
	fmt.Printf("%-18s", t.Format(time.Stamp))
	for _, coin := range *coins {
		price, err := strconv.ParseFloat(coin[fmt.Sprintf("price_%s", convert)], 64)
		if err != nil {
			log.Printf("couldn't convert the price to a float: %v", err)
			return
		}
		fmt.Printf("%-14g%-14g| ", price, price-(*prevPrices)[coin["id"]])
		(*prevPrices)[coin["id"]] = price
	}
	fmt.Println()
}
