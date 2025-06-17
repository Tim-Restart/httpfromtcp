package main

import (
	"fmt"
	"log"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/headers"
	"strings"
	"net/http"
	"io"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const host = 42069




// Handler functions go here

func firstHandler(w *response.Writer, req *request.Request) {

	inputRequest := req.RequestLine.RequestTarget
	// Checks for prefix of /httpbin and if true reassigns the inputRequest
	// This allows swtiching on that prefix

	if strings.HasPrefix(inputRequest, "/httpbin") {
		urlRequest := strings.TrimPrefix(inputRequest, "/httpbin")
		urlRequest = "https://httpbin.org" + urlRequest
		buf := make([]byte, 1024)

		res, err := http.Get(urlRequest)
		if err != nil {
			fmt.Println("Failed to make request")
		}
		defer res.Body.Close()

		err = w.WriteStatusLine(response.Ok) // Not sure if this is right
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/plain")
		headers.Set("Transfer-Encoding", "chunked")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}

		for {
			n, err := res.Body.Read(buf)
			if n > 0 {
				_, err2 := w.WriteChunkedBody(buf[:n])
				if err2 != nil {
					fmt.Println("Failed to write chunked body")
					return
				}
			}
			if err == io.EOF {
				break // end of stream
				}
			if err != nil {
				fmt.Println("Error doing chunked body")
				return
			}
		
		}

		_, err = w.WriteChunkedBodyDone()
			if err != nil {
				fmt.Println("Failed to end chunked body")
				return
			}
	}

	switch inputRequest {
	case "/yourproblem":
		err := w.WriteStatusLine(response.BadRequest) // Not sure if this is right
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/html")
		headers.Set("Connection", "close")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
		htmlContent := response.HtmlResponse(response.BadRequest)
		_, err = w.WriteBody([]byte(htmlContent))
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
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/html")
		headers.Set("Connection", "close")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
	
		htmlContent := response.HtmlResponse(response.InternalError)
		_, err = w.WriteBody([]byte(htmlContent))
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
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/html")
		headers.Set("Connection", "close")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
	
		htmlContent := response.HtmlResponse(response.Ok)
		_, err = w.WriteBody([]byte(htmlContent))
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

	




