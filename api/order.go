package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type orderType string
type orderDirection string

//Constants used for order creation.
const (
	Limit             orderType = "limit"
	Market            orderType = "market"
	FillOrKill        orderType = "fill-or-kill"
	ImmediateOrCancel orderType = "immediate-or-cancel"

	Buy  orderDirection = "buy"
	Sell orderDirection = "sell"
)

//The Fill struct represents a (partial) fulfillment of an order.
type Fill struct {
	Price    int       `json:"price"`
	Quantity int       `json:"qty"`
	TS       time.Time `json:"ts"`
}

//The Order struct contains information about an order.
type Order struct {
	ErrorResult
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
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders", baseURL, i.venue, i.symbol)
	i.RUnlock()

	if !i.setErr(jsonErr) {
		i.doHTTP("POST", url, bytes.NewBuffer(b), &v)
	}

	return
}

//CancelOrder cancels an order given it's id.
//See https://starfighter.readme.io/docs/cancel-an-order for further info about the actual API call.
func (i *Instance) CancelOrder(ID int) (v Order) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders/%s", baseURL, i.venue, i.symbol, strconv.Itoa(ID))
	i.RUnlock()

	i.doHTTP("DELETE", url, nil, &v)
	return
}

//OrderStatus returns the current order status for the given order id.
//See https://starfighter.readme.io/docs/status-for-an-existing-order for further info about the actual API call.
func (i *Instance) OrderStatus(ID int) (v Order) {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks/%s/orders/%s", baseURL, i.venue, i.symbol, strconv.Itoa(ID))
	i.RUnlock()

	i.doHTTP("GET", url, nil, &v)
	return
}

//AccountOrderStatus returns the current status for all orders of the current account on the current venue.
//See https://starfighter.readme.io/docs/status-for-all-orders for further info about the actual API call.
func (i *Instance) AccountOrderStatus() []Order {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/accounts/%s/orders", baseURL, i.venue, i.account)
	i.RUnlock()

	var v allOrdersStatusResult
	i.doHTTP("GET", url, nil, &v)
	return v.Orders
}

//StockOrderStatus returns the current status for all orders of the current stock on the current venue and account.
//See https://starfighter.readme.io/docs/status-for-all-orders-in-a-stock for further info about the actual API call.
func (i *Instance) StockOrderStatus() []Order {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/accounts/%s/stocks/%s/orders", baseURL, i.venue, i.account, i.symbol)
	i.RUnlock()

	var v allOrdersStatusResult
	i.doHTTP("GET", url, nil, &v)
	return v.Orders
}
