package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/herrberk/go-http2-streaming/http2"
)

var CLIENTS = 10
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
		client := new(http2.Client)
		client.Dial()

		// We just use this to start to print the info we want when we reach the last results.
		totalAmountOfReq := (CLIENTS * REQUESTS_PER_CLIENT * 95) / 100
		for i := 0; i < CLIENTS; i++ {
			go func(i int) {
				for j := 0; j < REQUESTS_PER_CLIENT; j++ {

					// Simulate processing time
					wait := 40 + rand.Intn(21)
					time.Sleep(time.Duration(wait) * time.Millisecond)

					go func() {
						startReqTime := time.Now()

						client.Post(data, REQUESTS_PER_CLIENT)
						reqTime := time.Since(startReqTime)
						totalRequestsDuration += reqTime
						done++

						if done > totalAmountOfReq {
							totalTime := time.Since(startTime)
							fmt.Println("total time: ", totalTime.Seconds())
							fmt.Println("avg request time per client: ", (totalRequestsDuration / time.Duration(CLIENTS)).Seconds())
						}
					}()
				}
			}(i)
			fmt.Println("done: " + strconv.Itoa(done))
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
