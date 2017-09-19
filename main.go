package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sharno/cryptoPrices/coinmarketcap"
)

func main() {

	limit := flag.Uint("limit", 2, "How many of the top coins on CoinMarketCap to be shown")
	convert := flag.String("convert", "USD",
		fmt.Sprintf("Convert the prices to a different currency other than USD, possible currencies are: %v", coinmarketcap.PossibleConvert))
	coins := flag.String("coins", "", "List the coins you want to monitor for prices separated by commas without spaces, Example: BTC,ETH")
	flag.Parse()

	convertUpperCase := strings.ToUpper(*convert)
	if !coinmarketcap.IsValidConvert(convertUpperCase) {
		log.Fatalf("The currency you entered couldn't be converted to: %s, please enter a valid currency", convertUpperCase)
	}

	if *coins != "" {
		coinsParsed := []string{}
		for _, v := range strings.Split(*coins, ",") {
			coinsParsed = append(coinsParsed, strings.ToUpper(v))
		}
		PrintPricesForeverOfCoins(coinsParsed, convertUpperCase)
	}

	PrintPricesForever(*limit, convertUpperCase)
}

// PrintPricesForever gets the prices from CoinMarketCap and prints them in the console every 5 minutes
func PrintPricesForever(limit uint, convert string) {
	fmt.Printf("CoinMarketCap changes the price every 5 minutes. The new price in %s would be printed here every 5 minutes\n\n", convert)

	prevPrices := &map[string]float64{}
	data, err := coinmarketcap.GetPrices(limit, convert)
	if err != nil {
		log.Printf("couldn't get the prices of coins: %v", err)
		return
	}
	PrintPrices(data, convert, time.Now(), true, prevPrices)

	ticker := time.Tick(5 * time.Minute)
	for t := range ticker {
		data, err := coinmarketcap.GetPrices(limit, convert)
		if err != nil {
			log.Printf("couldn't get the prices of coins: %v", err)
			return
		}
		PrintPrices(data, convert, t, false, prevPrices)
	}
}

// PrintPricesForeverOfCoins gets the individual coins prices from CoinMarketCap and prints them to console every 5 minutes
func PrintPricesForeverOfCoins(coins []string, convert string) {
	fmt.Printf("CoinMarketCap changes the price every 5 minutes. The new price in %s would be printed here every 5 minutes\n\n", convert)

	coinsIds := map[string]string{}
	allCoinsData, err := coinmarketcap.GetPrices(0, "USD")
	if err != nil {
		log.Printf("couldn't get prices of coins: %v", err)
	}
	for _, c := range *allCoinsData {
		coinsIds[c["symbol"]] = c["id"]
	}

	prevPrices := &map[string]float64{}
	data := &[]map[string]string{}
	for _, c := range coins {
		coin, err := coinmarketcap.GetCoinPrice(coinsIds[c], convert)
		if err != nil {
			log.Printf("couldn't get the price of coin %s: %v", c, err)
			return
		}
		*data = append(*data, *coin...)
	}
	PrintPrices(data, convert, time.Now(), true, prevPrices)

	ticker := time.Tick(5 * time.Minute)
	for t := range ticker {
		data := &[]map[string]string{}
		for _, c := range coins {
			coin, err := coinmarketcap.GetCoinPrice(coinsIds[c], convert)
			if err != nil {
				log.Printf("couldn't get the price of coin %s: %v", c, err)
				return
			}
			*data = append(*data, *coin...)
		}
		PrintPrices(data, convert, t, false, prevPrices)
	}
}

// PrintPrices gets the prices from CoinMarketCap and prints them in the console every 5 minutes
func PrintPrices(data *[]map[string]string, convert string, t time.Time, withNames bool, prevPrices *map[string]float64) {
	if withNames {
		fmt.Printf("%-18s", "timestamp")
		for _, coin := range *data {
			fmt.Printf("%-14s%-14s| ", coin["name"], "diff")
		}
		fmt.Println()
	}
	fmt.Printf("%-18s", t.Format(time.Stamp))
	for _, coin := range *data {
		price, err := strconv.ParseFloat(coin[fmt.Sprintf("price_%s", strings.ToLower(convert))], 64)
		if err != nil {
			log.Printf("couldn't convert the price to a float: %v", err)
			return
		}
		fmt.Printf("%-14g%-14g| ", price, price-(*prevPrices)[coin["id"]])
		(*prevPrices)[coin["id"]] = price
	}
	fmt.Println()
}
