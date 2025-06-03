package server

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"sync/atomic"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/request"
	"io"
	"bytes"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct{
	// Struct stuff here
	ServerState		atomic.Bool
	listener		net.Listener
	hand 			Handler
}


type HandlerError struct{
	HandlerStatusCode 	int
	HandlerMessage 		string
}

func (he HandlerError) Error() string {
	message := fmt.Sprintf("Status Code: %v Error: %v", he.HandlerStatusCode, he.HandlerMessage)
	return message
}



func Serve(port int, handler Handler) (*Server, error) {
	serverPort := ":"
	
	//Called in the main package to start the server
	s := &Server{}
	s.hand = handler
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

	parsedRequest, err := request.RequestFromReader(conn)
	if err != nil {
		_ = WriteHandlerError(conn, &HandlerError{})
		return
	}
	var buf bytes.Buffer

	err = s.hand(&buf, parsedRequest)
	if err != nil {
		_ = WriteHandlerError(conn, &HandlerError{})
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
		_ = WriteHandlerError(conn, &HandlerError{})
		return
	}

}

func WriteHandlerError(w io.Writer, he *HandlerError) error {


	err := response.WriteStatusLine(w, response.StatusCode(he.HandlerStatusCode))
	if err != nil {
		log.Println("Error writing handler status code")
		return err
	}

	_, err = w.Write([]byte(he.HandlerMessage))
	if err != nil {
		log.Println("Error writing handler message")
		return err
	}
	return nil
}



