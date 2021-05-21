# go-vestaboard

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/mikehelmick/go-vestaboard?tab=doc)
[![Go](https://github.com/mikehelmick/go-chaff/workflows/Go/badge.svg?event=push)](https://github.com/mikehelmick/go-vestaboard/actions?query=workflow%3AGo)

An unofficial client for the Vestaboard API in go.

# Usage

## Create a new client

```
client := vestaboard.New(c.APIKey, c.Secret)
```

From there, use the client methods

* `Viewer` to get the information from the connected viewer
* `Subscriptions` to get the subscription information
* `SendText` to post a message with the default formatting
  