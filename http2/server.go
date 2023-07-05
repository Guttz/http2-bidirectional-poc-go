package http2

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	httpRouter "github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httpRouter.Router
}

func (s *Server) Initialize() error {
	s.router = httpRouter.New()
	s.router.POST("/", s.handler)

	//Creates the http server
	server := &http.Server{
		Handler: s.router,
	}

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}

	log.Println("HTTP server is listening..")
	return server.ServeTLS(listener, "./http2/certs/key.crt", "./http2/certs/key.key")
}

func (s *Server) handler(w http.ResponseWriter, req *http.Request, _ httpRouter.Params) {
	// We only accept HTTP/2!
	// (Normally it's quite common to accept HTTP/1.- and HTTP/2 together.)
	if req.ProtoMajor != 2 {
		log.Println("Not a HTTP/2 request, rejected!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set headers related to event streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()

	buf := make([]byte, 4*1024)

	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			randomSeconds := rand.Intn(4) + 3
			time.Sleep(time.Duration(randomSeconds) * time.Second)
			w.Write([]byte("Server: Pushed message from server! \n"))
			flusher.Flush()
		}
	}()

	for {
		select {
		case <-notify:
			log.Println("HTTP connection just closed.")
			return
		default:
			// Write to the ResponseWriter
			n, err := req.Body.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
			}

			if err != nil {
				if err == io.EOF {
					w.Header().Set("Status", "200 OK")
					req.Body.Close()
				}
				break
			}

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()
		}

		// Pause for a second to simulate some data processing.
		time.Sleep(time.Second)
	}

}
