package main

import (
	"fmt"
	"log"
	"io"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const host = 42069


// Handler functions go here

func firstHandler(w io.Writer, req *request.Request) *server.HandlerError {
	log.Println("###$$$ You have entered the handler$$$###")
	log.Printf("Req: %v", req.RequestLine.RequestTarget)
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{HandlerStatusCode: 400, HandlerMessage: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{HandlerStatusCode: 500, HandlerMessage: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
	
}



func main() {

	fmt.Println("-----### Connection Establishing")

	server, err := server.Serve(host, firstHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("-----### Connection Established")

	defer server.Close()

	for {
	
	}
	return
	
}

	




