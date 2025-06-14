package server

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"sync/atomic"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/request"
)


type Handler func(w *response.Writer, req *request.Request)

type Server struct{
	// Struct stuff here
	ServerState		atomic.Bool
	listener		net.Listener
	hand 			Handler
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
	
	fmt.Println("-----### Server Established")

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
		// _ = WriteHandlerError(conn, &HandlerError{}) // Not needed after refactor
		return
	}

	writer := response.NewWriter(conn)
	s.hand(writer, parsedRequest) 
}

	