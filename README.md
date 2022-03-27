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

# Examples

There are a nice set of demos in cmd/
To use these, set your API Key and Secret as environment variables.
export APIKEY=YOUR_API_KEY_GOES_HERE
export SECRET=YOUR_SECRET_GOES_HERE

## Hello World

Does what it says - writes 'Hello World' to your vestaboard.

## Clock

## Game of Life

## Subscriptions

## Test Pattern

## Viewer