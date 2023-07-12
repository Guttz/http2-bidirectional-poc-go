package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/herrberk/go-http2-streaming/jsonrpc"
	"github.com/sourcegraph/jsonrpc2"
)

var CLIENTS = 1000
var REQUESTS_PER_CLIENT = 1000

func main() {
	// Waitc is used to hold the main function
	// from returning before the client connects to the server.
	waitc := make(chan bool)

	// Reads a file into memory - contains an example complete lsp request body
	data, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Println(err)
		return
	}

	startTime := time.Now()
	totalRequestsDuration := time.Duration(0)
	done := 0

	go func() {

		// We just use this to start to print the info we want when we reach the last results.
		totalAmountOfReq := (CLIENTS * REQUESTS_PER_CLIENT * 98) / 100
		for i := 0; i < CLIENTS; i++ {
			client := jsonrpc.Client()
			go func(client *jsonrpc2.Conn) {
				for j := 0; j < REQUESTS_PER_CLIENT; j++ {

					// Simulate processing time
					wait := 40 + rand.Intn(21)
					time.Sleep(time.Duration(wait) * time.Millisecond)

					go func() {
						startReqTime := time.Now()

						var result interface{}
						err = client.Call(context.Background(), "SomeMethod", data, &result)
						if err != nil {
							log.Fatalf("Call failed: %v", err)
						}

						reqTime := time.Since(startReqTime)
						totalRequestsDuration += reqTime
						done++

						if done > totalAmountOfReq {
							totalTime := time.Since(startTime)
							fmt.Println("total time: ", totalTime.Seconds())
							fmt.Println("avg request time per client: ", totalRequestsDuration.Seconds()/float64(CLIENTS*REQUESTS_PER_CLIENT))
						}
					}()
				}
			}(client)
		}
	}()

	// HTTP2 Server
	jsonrpc.Server()

	// Waits forever
	<-waitc
}
