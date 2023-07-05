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
