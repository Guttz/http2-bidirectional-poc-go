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

	// Set headers related to event streaming
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Status", "200 OK")
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Println(err)
	}
	wait := 40 + rand.Intn(21)
	time.Sleep(time.Duration(wait) * time.Millisecond)

	w.Write(body)
	w.WriteHeader(http.StatusOK)
	req.Body.Close()
}
