package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/herrberk/go-http2-streaming/http2"
)

var CLIENTS = 1000
var REQUESTS_PER_CLIENT = 1000

func main() {
	// Waitc is used to hold the main function
	// from returning before the client connects to the server.
	waitc := make(chan bool)

	// Reads a file into memory
	data, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Println(err)
		return
	}

	startTime := time.Now()
	totalRequestsDuration := time.Duration(0)
	done := 0
	// HTTP2 Client
	go func() {
		client := new(http2.Client)
		client.Dial()

		for i := 0; i < CLIENTS; i++ {
			go func(i int) {
				startReqTime := time.Now()
				client.Post(data, REQUESTS_PER_CLIENT)
				reqTime := time.Since(startReqTime)
				totalRequestsDuration += reqTime
				done++
				fmt.Println("done: " + strconv.Itoa(done))
				if done == CLIENTS {
					totalTime := time.Since(startTime)
					totalRequestsDuration = totalRequestsDuration / time.Duration(CLIENTS)
					fmt.Println("total time: ", totalTime.Seconds())
					fmt.Println("avg request time per client: ", totalRequestsDuration.Seconds())
				}
			}(i)
		}
	}()

	// HTTP2 Server
	server := new(http2.Server)
	err = server.Initialize()
	if err != nil {
		log.Println(err)
		return
	}

	// Waits forever
	<-waitc
}
