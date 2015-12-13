package api

import (
	"fmt"
	"net/http"
	"time"

	// use gorilla/websocket instead of x/net/websocket: https://github.com/gorilla/websocket#gorilla-websocket-compared-with-other-packages
	"github.com/gorilla/websocket"
)

const (
	baseWSURL string = "wss://api.stockfighter.io/ob/api/ws/"
)

type wsQuote struct {
	Ok    bool  `json:"ok"`
	Quote Quote `json:"quote"`
}

//The Execution struct gets only returned by websocket based calls.
type Execution struct {
	Order            Order     `json:"order"`
	StandingID       int       `json:"standingId"`
	IncomingID       int       `json:"incomingId"`
	Price            int       `json:"price"`
	Filled           int       `json:"filled"`
	FilledAt         time.Time `json:"filledAt"`
	StandingComplete bool      `json:"standingComplete"`
	IncomingComplete bool      `json:"incomingComplete"`
}

//QuotesForVenue returns a channel which streams all quotes for all symbols on the current venue.
//See https://starfighter.readme.io/docs/quotes-ticker-tape-websocket for further info about API call.
func (i *Instance) QuotesForVenue() <-chan Quote {
	urlExtension := fmt.Sprintf("%s/venues/%s/tickertape", i.account, i.venue)
	return i.wsQuotes(urlExtension)
}

//QuotesForStock returns a channel which streams all quotes for the current symbol on the current venue.
//See https://starfighter.readme.io/docs/quotes-ticker-tape-websocket for further info about API call.
func (i *Instance) QuotesForStock() <-chan Quote {
	urlExtension := fmt.Sprintf("%s/venues/%s/tickertape/stocks/%s", i.account, i.venue, i.symbol)
	return i.wsQuotes(urlExtension)
}

func (i *Instance) wsQuotes(urlExtension string) <-chan Quote {
	ch := make(chan Quote)

	go func() {
		conn, _, connErr := websocket.DefaultDialer.Dial(baseWSURL+urlExtension, http.Header{})
		defer close(ch)

		if connErr != nil {
			i.setErr(connErr)
		} else {
			//reuse variables
			var v wsQuote
			var jsonErr error

			for {
				jsonErr = conn.ReadJSON(&v)
				fmt.Println(v)
				if jsonErr != nil {
					conn.Close()
					i.setErr(jsonErr)
					break
				}
				ch <- v.Quote
			}
		}
	}()

	return ch
}

//ExecutionsForVenue returns a channel which streams all executions concerning the current account & venue.
//See https://starfighter.readme.io/docs/executions-fills-websocket for further info about API call.
func (i *Instance) ExecutionsForVenue() <-chan Execution {
	urlExtension := fmt.Sprintf("%s/venues/%s/executions", i.account, i.venue)
	return i.wsExecutions(urlExtension)
}

//ExecutionsForStock returns a channel which streams all executions concerning the current account, venue and symbol.
//See https://starfighter.readme.io/docs/executions-fills-websocket for further info about API call.
func (i *Instance) ExecutionsForStock() <-chan Execution {
	urlExtension := fmt.Sprintf("%s/venues/%s/executions/stocks/%s", i.account, i.venue, i.symbol)
	return i.wsExecutions(urlExtension)
}

func (i *Instance) wsExecutions(urlExtension string) <-chan Execution {
	ch := make(chan Execution)

	go func() {
		conn, _, connErr := websocket.DefaultDialer.Dial(baseWSURL+urlExtension, http.Header{})
		defer func() {
			conn.Close()
			close(ch)
		}()

		if connErr != nil {
			i.setErr(connErr)
		} else {
			//reuse variables
			var v Execution
			var jsonErr error

			for {
				jsonErr = conn.ReadJSON(&v)
				if jsonErr != nil {
					i.setErr(jsonErr)
					break
				}
				ch <- v
			}
		}
	}()

	return ch
}
