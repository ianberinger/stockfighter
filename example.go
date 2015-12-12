package main

import (
	"fmt"

	"github.com/ianberinger/stockfighter/api"
)

func main() {
	//start a new API instance with test values
	i := api.NewTestInstance()

	//see if the API is up
	fmt.Println("API is up:", i.Heartbeat())

	//see if the venue is up
	fmt.Println("Venue is up:", i.VenueHeartbeat())

	//available symbols at the current venue
	prettyPrint(i.AvailableStocks())

	//getting the current orderbook
	prettyPrint(i.Orderbook())
}

func prettyPrint(v interface{}) {
	fmt.Printf("%+v\n", v)
}
