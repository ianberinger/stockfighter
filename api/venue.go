package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Stock struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type availableStocksResult struct {
	Ok      bool    `json:"ok"`
	Symbols []Stock `json:"symbols"`
}

//AvailableStocks returns the available stock on a venue.
//See https://starfighter.readme.io/docs/list-stocks-on-venue for further info about the actual API call.
func (i *Instance) AvailableStocks() []Stock {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks", baseURL, i.venue)
	req, _ := http.NewRequest("GET", url, nil)
	i.RUnlock()
	res, httpErr := i.c.Do(req)
	i.setErr(httpErr)

	dec := json.NewDecoder(res.Body)
	var jsonErr error

	// API returns a different result format if there's an error
	if res.StatusCode == 200 {
		var v availableStocksResult
		jsonErr = dec.Decode(&v)
		return v.Symbols
	}

	var v errorResult
	jsonErr = dec.Decode(&v)
	i.setErr(apiError(v.Error, res.Status))

	i.setErr(jsonErr)
	return nil
}

//VenueHeartbeat works like Heartbeat() but for the current venue.
//See https://starfighter.readme.io/docs/venue-healthcheck for further info about the actual API call.
func (i *Instance) VenueHeartbeat() bool {
	i.RLock()
	urlExtension := fmt.Sprintf("venues/%s/heartbeat", i.venue)
	i.RUnlock()
	return i.heartbeat(urlExtension)
}
