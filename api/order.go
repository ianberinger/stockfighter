package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type orderType string
type orderDirection string

const (
	Limit             orderType = "limit"
	Market            orderType = "market"
	FillOrKill        orderType = "fill-or-kill"
	ImmediateOrCancel orderType = "immediate-or-cancel"

	Buy  orderDirection = "buy"
	Sell orderDirection = "sell"
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

type Fill struct {
	Price    int       `json:"price"`
	Quantity int       `json:"qty"`
	TS       time.Time `json:"ts"`
}

type Order struct {
	Account          string         `json:"account"`
	Venue            string         `json:"venue"`
	Symbol           string         `json:"symbol"`
	Price            int            `json:"price"`
	OriginalQuantity int            `json:"orignialQty"`
	Quantity         int            `json:"qty"`
	Direction        orderDirection `json:"direction"`
	OrderType        orderType      `json:"orderType"`
	ID               int            `json:"id"`
	TS               time.Time      `json:"ts"`
	Fills            []Fill         `json:"fills"`
	TotalFilled      int            `json:"totalFilled"`
	Open             bool           `json:"open"`
}

//NewOrder makes a new order and submits it to the API. See the package constants for available orderDirection and orderType types.
//NewOrder returns a Order struct of the created order.
//See https://starfighter.readme.io/docs/place-new-order for further info about the actual API call.
func (i *Instance) NewOrder(price int, quantity int, direction orderDirection, orderType orderType) (v Order) {
	i.RLock()
	b, jsonErr := json.Marshal(orderRequest{i.account, i.venue, i.symbol, price, quantity, direction, orderType})
	i.setErr(jsonErr)
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders", baseURL, i.venue, i.symbol)
	i.RUnlock()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header = i.h

	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
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

//CancelOrder cancels an order given it's id.
//See https://starfighter.readme.io/docs/cancel-an-order for further info about the actual API call.
func (i *Instance) CancelOrder(ID int) (v Order) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders/%s", baseURL, i.venue, i.symbol, strconv.Itoa(ID))
	req, _ := http.NewRequest("DELETE", url, nil)
	i.RUnlock()
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

//OrderStatus returns the current order status for the given order id.
//See https://starfighter.readme.io/docs/status-for-an-existing-order for further info about the actual API call.
func (i *Instance) OrderStatus(ID int) (v Order) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders/%s", baseURL, i.venue, i.symbol, strconv.Itoa(ID))
	req, _ := http.NewRequest("GET", url, nil)
	i.RUnlock()
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

type allOrdersStatusResult struct {
	Ok     bool    `json:"ok"`
	Venue  string  `json:"venue"`
	Orders []Order `json:"orders"`
}

//AccountOrderStatus returns the current status for all orders of the current account on the current venue.
//See https://starfighter.readme.io/docs/status-for-all-orders for further info about the actual API call.
func (i *Instance) AccountOrderStatus() []Order {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/accounts/%s/orders", baseURL, i.venue, i.account)
	req, _ := http.NewRequest("GET", url, nil)
	i.RUnlock()
	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
	var jsonErr error

	if res.StatusCode == 200 {
		var v allOrdersStatusResult
		jsonErr = dec.Decode(&v)
		return v.Orders
	}

	var v errorResult
	jsonErr = dec.Decode(&v)
	i.setErr(apiError(v.Error, res.Status))

	i.setErr(jsonErr)
	return nil
}

//StockOrderStatus returns the current status for all orders of the current stock on the current venue and account.
//See https://starfighter.readme.io/docs/status-for-all-orders-in-a-stock for further info about the actual API call.
func (i *Instance) StockOrderStatus() []Order {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/accounts/%s/stocks/%s/orders", baseURL, i.venue, i.account, i.symbol)
	req, _ := http.NewRequest("GET", url, nil)
	i.RUnlock()
	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
	var jsonErr error

	if res.StatusCode == 200 {
		var v allOrdersStatusResult
		jsonErr = dec.Decode(&v)
		return v.Orders
	}

	var v errorResult
	jsonErr = dec.Decode(&v)
	i.setErr(apiError(v.Error, res.Status))

	i.setErr(jsonErr)
	return nil
}
