package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//A MarketRequest struct represents an open position (either bid or ask() in the orderbook.
type MarketRequest struct {
	Price    int  `json:"price"`
	Quantity int  `json:"qty"`
	IsBuy    bool `json:"isBuy"`
}

//The Orderbook struct contains everything that gets returned on the Orderbook() API call.
type Orderbook struct {
	Venue  string          `json:"venue"`
	Symbol string          `json:"symbol"`
	Bids   []MarketRequest `json:"bids"`
	Asks   []MarketRequest `json:"asks"`
	TS     time.Time       `json:"ts"`
}

//Orderbook returns the orderbook for the current stock on the current venue.
//See: https://starfighter.readme.io/docs/get-orderbook-for-stock for further info about the actual API call.
//Returns an empty Oderbook struct if there was an error.
func (i *Instance) Orderbook() (v Orderbook) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s", baseURL, i.venue, i.symbol)
	req, _ := http.NewRequest("GET", url, nil)
	i.RUnlock()
	req.Header = i.h
	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
	var jsonErr error
	if res.StatusCode == 200 {
		jsonErr = dec.Decode(&v)
	} else {
		var v errorResult
		jsonErr = dec.Decode(&v)
		i.setErr(apiError(v.Error, res.Status))
	}

	i.setErr(jsonErr)
	return
}
