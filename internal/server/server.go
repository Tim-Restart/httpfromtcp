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
		_ = WriteHandlerError(conn, &HandlerError{})
		return
	}
	var buf bytes.Buffer

	err = s.hand(&buf, parsedRequest)
	if err != nil {
		value, ok := err.(*HandlerError)
		if ok {
			if value != nil {
				_ = WriteHandlerError(conn, value)
			}
		} else {
			unexpectedErr := &HandlerError{
				HandlerStatusCode: 500,
				HandlerMessage: "Server Error",
			}
		_ = WriteHandlerError(conn, unexpectedErr)
		log.Printf("Handler returned unexpected error type: %v", err)
		}
		return
	}

	err = response.WriteStatusLine(conn, response.Ok)
	if err != nil {
		log.Println("Error writing status code")
		return
	}
	length := buf.Len()
	defaultHeaders := response.GetDefaultHeaders(length)
	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		log.Println("Error writting headers")
		return
	}
	log.Printf("Buf print: %v", buf)
	_, err = buf.WriteTo(conn)
	if err != nil {
		_ = WriteHandlerError(conn, err.(*HandlerError))
		return
	}

}

func WriteHandlerError(w io.Writer, he *HandlerError) error {
	header := "Content-Type: text/plain\r\nContent-Length: "
	length := len(he.HandlerMessage)
	CRLF := "\r\n"
	headerFormated := []byte(header + strconv.Itoa(length) + CRLF + CRLF)

	// Creates the Status Line and sends it to the Writer
	err := response.WriteStatusLine(w, response.StatusCode(he.HandlerStatusCode))
	//log.Printf("##--WriteHandlerError Response: %v", response.StatusCode(he.HandlerStatusCode))
	if err != nil {
		log.Println("Error writing handler status code")
		return err
	}

	// Create the headers and send to the writer
	_, err = w.Write(headerFormated)
	if err != nil {
		log.Println("Error writing headers")
		return err
	}

	_, err = w.Write([]byte(he.HandlerMessage))
	//log.Printf("##--WriteHanderError Write: %v", he.HandlerMessage)
	if err != nil {
		log.Println("Error writing handler message")
		return err
	}
	return nil
}

/*

func WriteSuccessM(w io.Writer, message string) error {

	length := len(message)

	err = response.WriteStatusLine(w, response.Ok)
	if err != nil {
		log.Println("Error writing status code")
		return err
	}
	length := buf.Len()
	defaultHeaders := response.GetDefaultHeaders(length)
	err = response.WriteHeaders(w, defaultHeaders)
	if err != nil {
		log.Println("Error writting headers")
		return err
	}
	log.Printf("Buf print: %v", buf)
	_, err = buf.WriteTo(conn)
	if err != nil {
		_ = WriteHandlerError(conn, err.(*HandlerError))
		return err
	}
	return nil
}

*/
