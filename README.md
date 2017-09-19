Crypto currencies prices monitor in Go
======================================


I like to watch the change of prices of crypto currencies. A terminal window with such a monitor is much lighter on system resources. That's why I made this command line tool.

It's still in the beginning and I'm trying to make it a comprehensive tool for such a thing

The tool just connects now to CoinMarketCap API. Later I'll add more APIS to connect to. The API updates the prices every 5 minutes.

usage flags:

`-limit <positive number>`  gets the highest number of cryptocurrencies in the market (default: 2)

`-convert <fiat currency symbol>` gets the prices of the cryptocurrencies in a specific fiat currency like USD, EUR ... etc (default: USD)

`-coins <list of cryptocurrency symbols separated by commas>` ignores the limit and gets the list of cryptocurrencies specified by this option (example: BTC,ETH,XMR)
