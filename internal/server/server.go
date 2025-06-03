package server

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"sync/atomic"
	"httpfromtcp/internal/response"
	"io"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct{
	// Struct stuff here
	ServerState		atomic.Bool
	listener		net.Listener
	handle 			Handler
}


type HandlerError struct{
	handlerStatusCode 	int
	handlerMessage 		string
}





func Serve(port int, handler Handler) (*Server, error) {
	serverPort := ":"
	
	//Called in the main package to start the server
	s := &Server{}
	s.handle = handler
	serverPort += strconv.Itoa(port)
	
	s.ServerState.Store(true)
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener

	go s.listen()
	
	fmt.Println("--##--Server Established--##--")

	return s, nil
}



func (s *Server) Close() error {
	// Closes the listener and the server
	err := s.listener.Close()
	s.ServerState.Store(false)
	return err
}


func (s *Server) listen() {
	// uses a loop to .Accept new connections as they come in, and handles
	// each one in a new goroutine. 
	// Uses atomic.Bool to track whether the server is closed
	// or not so that it can ignore connection errors after it is closed
	for {
		
		serverConnection := s.ServerState.Load()
		if serverConnection != true {
			
			break
		}
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go s.handle(conn)
		
		// Send it off to the handler here?
		}
	}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	parsedRequest, err := RequestFromReader(conn)
	if err != nil {
		_ := WriteHandlerError(conn, err)
		return
	}
	var buf bytes.Buffer

	err = s.handler(buf, parsedRequest)
	if err != nil {
		_ := WriteHandlerError(conn, err)
		return
	}

	err = response.WriteStatusLine(conn, response.Ok)
	if err != nil {
		log.Println("Error writing status code")
		return
	}
	defaultHeaders := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		log.Println("Error writting headers")
		return
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		_ := WriteHandlerError(conn, err)
		return
	}

}

func WriteHandlerError(w io.Writer, he *HandlerError) error {


	err := response.WriteStatusLine(w, he.handlerStatusCode)
	if err != nil {
		log.Println("Error writing handler status code")
		return err
	}

	_, err := w.Write([]byte(he.handlerMessage))
	if err != nil {
		log.Println("Error writing handler message")
		return err
	}

}



