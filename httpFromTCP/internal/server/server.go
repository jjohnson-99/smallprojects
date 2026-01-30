package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/jjohnson-99/httpFromTCP/internal/request"
	"github.com/jjohnson-99/httpFromTCP/internal/response"
)

type Server struct {
	handler    Handler
	inShutdown atomic.Bool
	listener   net.Listener
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message	   string
}

type Handler func(w *response.Writer, req *request.Request) //*HandlerError

func Serve(port int, handler Handler) (*Server, error) {	
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) 
	if err != nil {
		return nil, err
	}
	server := &Server{
		handler:  handler,
		listener: listener,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.inShutdown.Store(true)
	s.listener.Close()

	return nil
}

func (s *Server) listen() {
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.inShutdown.Load() {
				return
			}
			log.Fatal(err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message: 	err.Error(),
		}
		handlerError.Write(conn)
		return
	}

	w := &response.Writer{Buffer: bytes.NewBuffer([]byte{})}
	s.handler(w, req)
	b := w.Buffer.Bytes()
	conn.Write(b)
	fmt.Println(string(w.Buffer.Bytes()))
	return
}



//func httpbinHandler(w *response.Writer, req *request.Request) {
	//resp, err := net.http.Get(url)

//}

func (e *HandlerError) Write(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%d: %s", e.StatusCode, e.Message)))
}
