package http2

import (
	"bufio"
	"math/rand"

	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

type Client struct {
	client *http.Client
}

func (c *Client) Dial() {
	// Adds TLS cert-key pair
	certs, err := tls.LoadX509KeyPair("./http2/certs/key.crt", "./http2/certs/key.key")
	if err != nil {
		log.Fatal(err)
	}

	t := &http2.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{certs},
			InsecureSkipVerify: true,
		},
	}

	c.client = &http.Client{Transport: t}
}

func (c *Client) Post(data []byte, requestsPerClient int) {
	// Create a pipe to read and write data
	pr, pw := io.Pipe()

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "localhost:8080",
			Path:   "/",
		},
		Header: http.Header{},
		Body:   pr,
	}

	// In a separate goroutine, write data to the request body
	go func() {
		// We don't close so we keep this socket alive
		//defer pw.Close()

		// Write to the pipe
		for i := 0; i < 1; i++ {
			pw.Write([]byte("server intialized"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Sends the request
	resp, err := c.client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode == 500 {
		return
	}

	defer resp.Body.Close()

	bufferedReader := bufio.NewReader(resp.Body)

	buffer := make([]byte, 4*1024)

	var totalBytesReceived int

	reachedEOF := make(chan bool)
	// Reads the response
	go func() {
		for {
			len, err := bufferedReader.Read(buffer)

			if len > 0 {
				totalBytesReceived += len
				//fmt.Println("\nEchoed msg: " + string(buffer[:len]))
			}

			if err != nil {
				if err == io.EOF {
					// Last chunk received
					log.Println(err)
					reachedEOF <- true
				}
				break
			}
		}
	}()

	for i := 0; i < requestsPerClient; i++ {
		/* time.Sleep(300 * time.Millisecond)
		fmt.Print("Msg to send: ")
		stdin := bufio.NewReader(os.Stdin)
		msg, _ := stdin.ReadString('\n') */
		//pw.Write([]byte(data))
		//pw.Write([]byte("hi"))
		//time.Sleep(100 * time.Millisecond)

		// Simulate user typing
		wait := 40 + rand.Intn(21)
		time.Sleep(time.Duration(wait) * time.Millisecond)

		pw.Write(data)
	}

	pw.Close()
	<-reachedEOF
}
