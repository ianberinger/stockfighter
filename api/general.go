// Package api contains a wrapper for the stockfighter.io API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const (
	baseURL string = "https://api.stockfighter.io/ob/api/"
)

type err struct {
	sync.Mutex
	v error
}

//Instance is the basic unit of operation for all API actions.
type Instance struct {
	//not protected by mutex because they don't get touched by us.
	c http.Client
	h http.Header

	//protected by it's own mutex
	err err

	sync.RWMutex
	account string
	venue   string
	symbol  string
}

//setErr sets the error value on an instance only when the error isn't nil.
//This is useful because we don't have to check if the error is nil before calling setErr() if we don't want a current error overwriten by a nil one.
func (i *Instance) setErr(err error) {
	if err != nil {
		i.err.Lock()
		i.err.v = err
		i.err.Unlock()
	}
}

//Err returns the last error of an API instance.
func (i *Instance) Err() error {
	i.err.Lock()
	defer i.err.Unlock()
	return i.err.v
}

//ResetErr resets the error of an API instance to nil.
func (i *Instance) ResetErr() {
	i.err.Lock()
	i.err.v = nil
	i.err.Unlock()
}

//GetAccount gets the current account of an instance.
func (i *Instance) GetAccount() string {
	i.RLock()
	defer i.RUnlock()
	return i.account
}

//GetVenue gets the current venue of an instance.
func (i *Instance) GetVenue() string {
	i.RLock()
	defer i.RUnlock()
	return i.venue
}

//GetSymbol gets the current stock symbol of an instance.
func (i *Instance) GetSymbol() string {
	i.RLock()
	defer i.RUnlock()
	return i.symbol
}

//SetAccount changes the current account of an instance. Waits until all current read operations are completed and blocks while changing.
func (i *Instance) SetAccount(account string) {
	i.Lock()
	i.account = account
	i.Unlock()
}

//SetVenue changes the current venue of an instance. Waits until all current read operations are completed and blocks while changing.
func (i *Instance) SetVenue(venue string) {
	i.Lock()
	i.venue = venue
	i.Unlock()
}

//SetSymbol changes the current stock symbol of an instance. Waits until all current read operations are completed and blocks while changing.
func (i *Instance) SetSymbol(symbol string) {
	i.Lock()
	i.symbol = symbol
	i.Unlock()
}

//NewInstance creates a new API instance based on the given inputs and returns a pointer to it.
func NewInstance(APIKey, account, venue, symbol string) *Instance {
	// create default header
	h := http.Header{}
	h.Add("X-Starfighter-Authorization", APIKey)

	return &Instance{http.Client{}, h, err{}, sync.RWMutex{}, account, venue, symbol}
}

//NewTestInstance calls NewInstance with useful presets for package testing.
func NewTestInstance() *Instance {
	return NewInstance("", "EXB123456", "TESTEX", "FOOBAR")
}

//Heartbeat checks if the API is up and returns true if it is.
//See: https://starfighter.readme.io/docs/heartbeat for further info about API call.
func (i *Instance) Heartbeat() bool {
	return i.heartbeat("heartbeat")
}

type errorResult struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

//apiError creates an error from an error message and a http status.
func apiError(str string, status string) error {
	if str != "" {
		return fmt.Errorf("API: %s; %s", status, str)
	}
	return nil
}

func (i *Instance) heartbeat(urlExtension string) bool {
	req, _ := http.NewRequest("GET", baseURL+urlExtension, nil)
	req.Header = i.h
	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
	var v errorResult
	jsonErr := dec.Decode(&v)
	i.setErr(jsonErr)

	if httpErr == nil && jsonErr == nil {
		if !v.Ok {
			i.setErr(apiError(v.Error, res.Status))
		}
		return v.Ok
	}
	return false
}
