package main

import (
	"fmt"
	"log"
	"io"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"httpfromtcp/internal/response"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const host = 42069

func NewWriter( w io.Writer) *Writer{
	newWriter := &Writer{
		writer: w,
		WriteStatus: 0,
	}
	return newWriter
}


// Handler functions go here

func firstHandler(w *response.Writer, req *request.Request) {

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		err := w.WriteStatusLine(response.BadRequest) // Not sure if this is right
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
		err = w.WriteBody([]byte(htmlContent))
		if err != nil {
			fmt.Println("Failed to write body")
			return
		}
	
	case "/myproblem":
		err := w.WriteStatusLine(response.InternalError) // Not sure if this is right
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
		err = w.WriteBody([]byte(htmlContent))
		if err != nil {
			fmt.Println("Failed to write body")
			return
		}
	default:
			err := w.WriteStatusLine(response.Ok) // Not sure if this is right
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
		err = w.WriteBody([]byte(htmlContent))
		if err != nil {
			fmt.Println("Failed to write body")
			return
		}
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

	




