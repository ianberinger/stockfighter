package api

import (
	"fmt"
	"net/http"
	"time"

	// use gorilla/websocket instead of x/net/websocket: https://github.com/gorilla/websocket#gorilla-websocket-compared-with-other-packages
	"github.com/gorilla/websocket"
)

var (
	baseWSURL string = "wss://api.stockfighter.io/ob/api/ws/"
)

type wsQuote struct {
	ErrorResult
	Quote Quote `json:"quote"`
}

type streamer interface {
	Stop()
	Stopped() bool
	add(interface{})
	close()
}

//QuoteStream contains a Values chan which streams Quotes. Implements streamer interface.
type QuoteStream struct {
	Values chan Quote
	stop   bool
}

func SetBaseWSURL(URL string) {
	baseWSURL = URL
}

//Stop stops the Stream.
func (s *QuoteStream) Stop() {
	s.stop = true
}

//Stopped returns true if the stream was stopped.
func (s *QuoteStream) Stopped() bool {
	return s.stop
}

func (s *QuoteStream) add(v interface{}) {
	s.Values <- v.(*wsQuote).Quote
}

func (s *QuoteStream) close() {
	close(s.Values)
}

//ExecutionStream contains a Values chan which streams Executions. Implements streamer interface.
type ExecutionStream struct {
	Values chan Execution
	stop   bool
}

//Stop stops the Stream.
func (s *ExecutionStream) Stop() {
	s.stop = true
}

//Stopped returns true if the stream was stopped.
func (s *ExecutionStream) Stopped() bool {
	return s.stop
}

func (s *ExecutionStream) add(v interface{}) {
	s.Values <- *v.(*Execution)
}

func (s *ExecutionStream) close() {
	close(s.Values)
}

//The Execution struct gets only returned by websocket based calls.
type Execution struct {
	ErrorResult
	Order            Order     `json:"order"`
	StandingID       int       `json:"standingId"`
	IncomingID       int       `json:"incomingId"`
	Price            int       `json:"price"`
	Filled           int       `json:"filled"`
	FilledAt         time.Time `json:"filledAt"`
	StandingComplete bool      `json:"standingComplete"`
	IncomingComplete bool      `json:"incomingComplete"`
}

func (i *Instance) wsURL(method string, stockOnly bool, account string) string {
	if stockOnly {
		return fmt.Sprintf("%s%s/venues/%s/%s/stocks/%s", baseWSURL, account, i.venue, method, i.symbol)
	}
	return fmt.Sprintf("%s%s/venues/%s/%s", baseWSURL, account, i.venue, method)
}

//Quotes returns a stream which streams all quotes for the current venue or only the current stock.
//A stream can be terminated with: stream.Stop()
//See https://starfighter.readme.io/docs/quotes-ticker-tape-websocket for further info about API call.
func (i *Instance) Quotes(stockOnly bool) *QuoteStream {
	s := &QuoteStream{make(chan Quote), false}
	go i.doWS(s, i.wsURL("tickertape", stockOnly, i.account), &wsQuote{})
	return s
}

//Executions returns a stream which streams all executions for the current venue or only the current stock.
//Authentication is done with the account number
//A stream can be terminated with: stream.Stop()
//See https://starfighter.readme.io/docs/executions-fills-websocket for further info about API call.
func (i *Instance) Executions(stockOnly bool, account string) *ExecutionStream {
	s := &ExecutionStream{make(chan Execution), false}
	go i.doWS(s, i.wsURL("executions", stockOnly, account), &Execution{})
	return s
}

func (i *Instance) doWS(s streamer, url string, v apiResponse) {
	conn, _, connErr := websocket.DefaultDialer.Dial(url, http.Header{})
	defer func() {
		conn.Close()
		s.close()
	}()

	if !i.setErr(connErr) {
		for !i.setErr(conn.ReadJSON(v)) {
			if v.isOk() && !s.Stopped() {
				s.add(v)
			} else {
				i.setErr(v.err("WS"))
				s.Stop()
				break
			}
		}
	}
}
