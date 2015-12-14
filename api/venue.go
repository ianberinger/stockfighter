package api

import "fmt"

//The Stock struct contains the name and the ticker symbol of a stock.
type Stock struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

//AvailableStocks returns the available stock on a venue.
//See https://starfighter.readme.io/docs/list-stocks-on-venue for further info about the actual API call.
func (i *Instance) AvailableStocks() []Stock {
	i.RLock()
	url := fmt.Sprintf("%svenues/%s/stocks", baseURL, i.venue)
	i.RUnlock()

	var v availableStocksResult
	i.doHTTP("GET", url, nil, &v)

	return v.Symbols
}

//VenueHeartbeat works like Heartbeat() but for the current venue.
//See https://starfighter.readme.io/docs/venue-healthcheck for further info about the actual API call.
func (i *Instance) VenueHeartbeat() ErrorResult {
	i.RLock()
	urlExtension := fmt.Sprintf("venues/%s/heartbeat", i.venue)
	i.RUnlock()

	return i.heartbeat(urlExtension)
}
