// Package api contains a wrapper for the stockfighter.io API.
package api

import (
	"net/http"
	"sync"
)

//Instance is the basic unit of operation for all API actions.
type Instance struct {
	debug bool
	//not protected by mutex because they don't get touched by us.
	c http.Client
	h http.Header

	//each protected by it's own mutex
	err err
	state
}

type state struct {
	sync.RWMutex
	instanceID int
	account    string
	venue      string
	symbol     string
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

//SetAPIKey changes the API-Key of an instance.
func (i *Instance) SetAPIKey(apiKey string) {
	i.h.Set("X-Starfighter-Authorization", apiKey)
}

//SetInstanceID changes the current instanceID. Waits until all current read operations are completed and blocks while changing.
func (i *Instance) SetInstanceID(instanceID int) {
	i.Lock()
	i.instanceID = instanceID
	i.Unlock()
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

//setState sets whole state in one lock op.
func (i *Instance) setState(instanceID int, account, venue, symbol string) {
	i.Lock()
	i.instanceID = instanceID
	i.account = account
	i.venue = venue
	i.symbol = symbol
	i.Unlock()
}

//New creates a new API instance without any presets.
func New(apiKey string) (i *Instance) {
	i = &Instance{}
	i.c = http.Client{}
	i.h = http.Header{}
	i.SetAPIKey(apiKey)
	return
}

//NewInstance creates a new API instance and sets defaults.
//Shorcut for New() -> SetAccount() -> SetVenue() -> SetSymbol()
func NewInstance(apiKey, account, venue, symbol string) (i *Instance) {
	i = New(apiKey)
	i.setState(0, account, venue, symbol)
	return
}

//NewTestInstance calls NewInstance with useful presets for package testing.
func NewTestInstance() *Instance {
	return NewInstance("", "EXB123456", "TESTEX", "FOOBAR")
}

//Heartbeat checks if the API is up and returns true if it is.
//See: https://starfighter.readme.io/docs/heartbeat for further info about API call.
func (i *Instance) Heartbeat() ErrorResult {
	return i.heartbeat("heartbeat")
}

func (i *Instance) heartbeat(urlExtension string) (v ErrorResult) {
	i.doHTTP("GET", baseURL+urlExtension, nil, &v)
	return
}
