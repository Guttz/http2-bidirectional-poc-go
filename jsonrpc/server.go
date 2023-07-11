package jsonrpc

import (
	"context"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

func handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	switch req.Method {
	case "initialize":
		return nil, nil
		//return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	}
	return nil, nil
}

func Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	// AUTH - We could here do authentication to check if some token is valid everytime we receive a request, for example

	// Prevent any uncaught panics from taking the entire server down.
	defer func() {
		if perr := panicf(recover(), "%v", req.Method); perr != nil {
			err = perr
		}
	}()

	res, err := handle(ctx, conn, req)
	if err != nil {
		log.Printf("error serving, %+v\n", err)
	}

	return res, err
}

func Server() {
	lis, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("Failed to accept: %v", err)
		}

		stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})

		jsonrpc2.NewConn(
			context.Background(),
			stream,
			jsonrpc2.HandlerWithError(Handle),
			nil,
		)
	}
}
