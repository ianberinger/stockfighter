package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

var (
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
	ErrorResult
	Venue  string  `json:"venue"`
	Orders []Order `json:"orders"`
}

type availableStocksResult struct {
	ErrorResult
	Symbols []Stock `json:"symbols"`
}

type apiResponse interface {
	isOk() bool
	err(status string) error
}

//ErrorResult gets returned inside every response of an API method.
//Allows for an alternative way of error checking (without calling i.GetErr()).
type ErrorResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"error"`
}

func SetBaseURL(URL string) {
	baseURL = URL
}

func (e ErrorResult) isOk() bool {
	return e.Ok
}

func (e ErrorResult) err(status string) error {
	str := e.Message
	if str == "" {
		str = "no error message"
	}
	return fmt.Errorf("API: %s; %s", status, str)
}

func (i *Instance) doHTTP(httpVerb string, url string, body io.Reader, v apiResponse) {
	req, err := http.NewRequest(httpVerb, url, body)
	req.Header = i.h

	if i.debug && !i.setErr(err) {
		reqDump, err := httputil.DumpRequestOut(req, true)
		i.setErr(err)
		fmt.Printf("request: %s", reqDump)
	}

	res, err := i.c.Do(req)
	if i.debug && !i.setErr(err) {
		resDump, err := httputil.DumpResponse(res, true)
		i.setErr(err)
		fmt.Printf("response: %s", resDump)
	}

	i.setErr(json.NewDecoder(res.Body).Decode(v))

	if !v.isOk() || res.StatusCode != 200 {
		i.setErr(v.err(res.Status))
	}
}
