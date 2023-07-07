package http2

import (
	"bytes"
	"io"

	"crypto/tls"
	"log"
	"net/http"
	"net/url"

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

func nothing(data []byte) {

}

func (c *Client) Post(data []byte, requestsPerClient int) {
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "localhost:8080",
			Path:   "/",
		},
		Header: http.Header{},
		Body:   io.NopCloser(bytes.NewReader(data)),
	}

	// Sends the request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode == 500 {
		return
	}

	// we just want to make sure we read the body so it' similar to what we do when using a single request
	body, err := io.ReadAll(resp.Body)
	nothing(body)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		//log.Println("POST request succeeded")
	} else {
		//log.Printf("POST request failed with status code %d", resp.StatusCode)
	}
}
