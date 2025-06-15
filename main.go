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

func firstHandler(w *response.Writer, req *request.Request) *server.HandlerError {

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		preparedResponse := //somefunction call here to the response.Writer package
		return &server.HandlerError{HandlerStatusCode: 400, HandlerMessage: "Your problem is not my problem\n"}
	case "/myproblem":
		preparedResponse := //somefunction call here to the response.Writer package
		return &server.HandlerError{HandlerStatusCode: 500, HandlerMessage: "Woopsie, my bad\n"}
	default:
		
		preparedResponse := //somefunction call here to the response.Writer package
		return &server.HandlerError{HandlerStatusCode: 200, HandlerMessage: "All good, frfr\n"}
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

	




