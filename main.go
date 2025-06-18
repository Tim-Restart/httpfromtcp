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
	"crypto/sha256"
	"strconv"
	"os"
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
		header := headers.NewHeaders()
		header.Set("Content-Type", "text/plain")
		header.Set("Transfer-Encoding", "chunked")
		header.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		
		err = w.WriteHeaders(header)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}

		var fullBody []byte

		for {
			n, err := res.Body.Read(buf)
			if n > 0 {
				_, err2 := w.WriteChunkedBody(buf[:n])
				fullBody = append(fullBody, buf[:n]...)
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

		sum := fmt.Sprintf("%x", sha256.Sum256(fullBody))
		chunkedLength := strconv.Itoa(len(fullBody))
		trailers := headers.NewHeaders()
		trailers.Set("X-Content-SHA256", sum)
		trailers.Set("X-Content-Length", chunkedLength)
		err = w.WriteTrailers(trailers)
		if err != nil {
			fmt.Println("Failed to write trailers")
			return
		}
		return
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

	case "/video":
		err := w.WriteStatusLine(response.Ok)
		if err != nil {
			fmt.Println("Failed to write status line")
			return
		}
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "video/mp4")
		headers.Set("Connection", "close")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println("Failed to write headers")
			return
		}
		data, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.WriteBody(data)
		if err != nil {
			fmt.Println("Failed to watch video")
			return
		}
		return


	
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

	




