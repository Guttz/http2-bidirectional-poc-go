package main

import (
	"github.com/herrberk/go-http2-streaming/jsonrpc"
)

func main() {
	go jsonrpc.Server()
	jsonrpc.Client()
}
