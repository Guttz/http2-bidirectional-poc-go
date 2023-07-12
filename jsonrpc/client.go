package jsonrpc

import (
	"context"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

type SomeType struct {
	// Define your type here
}

type noopHandler struct{}

func (noopHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {}

func Client() *jsonrpc2.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})
	client := jsonrpc2.NewConn(
		context.Background(),
		stream,
		noopHandler{},
		nil,
	)

	return client

}
