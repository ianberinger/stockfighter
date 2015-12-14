package api

import (
	"fmt"
	"time"
)

//The Quote struct contains all data that gets returned on a Quote() call.
type Quote struct {
	ErrorResult
	Venue     string    `json:"venue"`
	Symbol    string    `json:"symbol"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	BidSize   int       `json:"bidSize"`
	AskSize   int       `json:"askSize"`
	BidDepth  int       `json:"bidDepth"`
	AskDepth  int       `json:"askDepth"`
	LastPrice int       `json:"last"`
	LastSize  int       `json:"lastSize"`
	LastTrade time.Time `json:"lastTrade"`
	QuoteTime time.Time `json:"quoteTime"`
}

//Quote returns the quote for the current stock on the current venue.
//See https://starfighter.readme.io/docs/a-quote-for-a-stock for further info about the actual API call.
//Returns an empty Oderbook struct if there was an error.
func (i *Instance) Quote() (v Quote) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s/quote", baseURL, i.venue, i.symbol)
	i.RUnlock()

	i.doHTTP("GET", url, nil, &v)
	return
}
