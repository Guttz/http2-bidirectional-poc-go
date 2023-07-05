package http2

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

func (c *Client) Post(data []byte) {
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
		// TODO - remove comment
		defer pw.Close()

		// Write to the pipe
		for i := 0; i < 5; i++ {
			pw.Write([]byte("streaming data " + strconv.Itoa(i) + "\n"))
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// Sends the request
	log.Println("Executing request..")
	resp, err := c.client.Do(req)
	log.Println("Request sent..")

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

	// Reads the response
	for {
		len, err := bufferedReader.Read(buffer)

		if len > 0 {
			totalBytesReceived += len
			log.Println("Client:\n" + string(buffer[:len]))
		}

		if err != nil {
			if err == io.EOF {
				// Last chunk received
				log.Println(err)
			}
			break
		}
	}

	/* 	fmt.Print("Msg to send: ")
	   	stdin := bufio.NewReader(os.Stdin)
	   	msg, err := stdin.ReadString('\n')
	   	pw.Write([]byte(msg)) */
	//log.Println("Total Bytes Sent:", len(data))
}
