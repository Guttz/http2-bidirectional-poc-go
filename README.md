## POC Bidirectional Streaming with http2

Install packages:

```sh
$ go mod tidy
```

Start the demo: 

```sh
$ go run main.go
```

You should send a prompt asking for the message to send. If you write a message and press enter, it's going to be sent to the server and the server will respond with the same message back. Sometimes the server will also push random messages to show it can also push messages to the client.

We basically have one http2 full duplex connection that can be used from both server and client to communicate to one another.

## Gotchas

HTTP2 will buffer the body as default and wait for it to be closed before sending it to the server and also sending the response back to the client. So on the server side we have to convert our writer to a http flusher and flush the buffer everytime we want to immediately send the buffer content to the client. On the client side we create a in-memory pipe and write. 

With HTTP2 TLS is required so we also have to configure the certificates.

The frame we receive on the client and in the server can ofc contain multiple messages that we're stacked before being flushed to the client. So this is one of the challenges we have to deal with. If we lock/sync everything we lose the perfomance. And if we allow sending multiple objects at a time we have to deserialize them properly on the client side and in the server.
