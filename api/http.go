package api

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	baseURL string = "https://api.stockfighter.io/ob/api/"
)

type orderRequest struct {
	Account   string         `json:"account"`
	Venue     string         `json:"venue"`
	Symbol    string         `json:"symbol"`
	Price     int            `json:"price"`
	Quantity  int            `json:"qty"`
	Direction orderDirection `json:"direction"`
	OrderType orderType      `json:"orderType"`
}

type allOrdersStatusResult struct {
	Ok     bool    `json:"ok"`
	Venue  string  `json:"venue"`
	Orders []Order `json:"orders"`
}

type availableStocksResult struct {
	Ok      bool    `json:"ok"`
	Symbols []Stock `json:"symbols"`
}

func (i *Instance) doHTTP(httpVerb string, url string, body io.Reader, v interface{}) {
	req, _ := http.NewRequest(httpVerb, url, body)
	req.Header = i.h

	res, err := i.c.Do(req)
	i.setErr(err)

	if res.StatusCode == 200 {
		err = json.NewDecoder(res.Body).Decode(v)
	} else {
		var v errorResult
		err = json.NewDecoder(res.Body).Decode(&v)
		i.setErr(apiError(v.Error, res.Status))
	}
	i.setErr(err)
}
