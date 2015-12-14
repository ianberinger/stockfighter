# stockfighter
[![Build Status](https://travis-ci.org/ianberinger/stockfighter.svg?branch=master)](https://travis-ci.org/ianberinger/stockfighter) [![GoDoc](https://godoc.org/github.com/ianberinger/stockfighter?status.svg)](https://godoc.org/github.com/ianberinger/stockfighter/api)

Stateful API client for stockfighter.io written in Go

Implements all Trade API calls from the documentation (https://starfighter.readme.io) and some of the GameMaster API calls.

Trade API calls are tested and work, GM API calls should mostly work.

### Usage
You can use the API client in your own Go packages:

	import "github.com/ianberinger/stockfighter/api"

Documentation is on [godoc](https://godoc.org/github.com/ianberinger/stockfighter/api) .
See [example.go](./example.go) for a usage example.
