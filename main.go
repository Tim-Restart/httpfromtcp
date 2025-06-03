package main

import (
	"fmt"
	"log"
	"io"
	"strings"
	"net"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const host = ":42069"

func getLinesChannel(conn net.Conn) <-chan string {

	ch := make(chan string)
	
	go func() {
		defer close(ch)
		buff := make([]byte, 8)
		var currentLine string
		
		for {
		n, err := conn.Read(buff) // This line needs to be changed to take the incomming connection
		
		if err == io.EOF {  // This error needs to be changed to close when the connection closes or similar
			// More to read
			if currentLine != "" {
				ch <- currentLine
			}
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

			// Process the entire chunk of bytes read
			content := string(buff[:n])
			for len(content) > 0 {
				// Find the next newline in the content
				i := strings.Index(content, "\n")
				if i >= 0 {
					// We found a newline
					currentLine += content[:i] // Add everything up to the newline
					ch <- currentLine          // Send the complete line
					currentLine = ""           // Reset for the next line
					content = content[i+1:]    // Skip past the newline
				} else {
					// No more newlines in this chunk
					currentLine += content
					break
				}
				}

			}
		}()
	return ch	
}

// Handler functions go here

func firstHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{handlerStatusCode: 400, handlerMessage: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{handlerStatusCode: 500, handlerMessage: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
	
}



func main() {

	fmt.Println("-----### Connection Establishing ###-----")

	listen, err := net.Listen(network, host)
		if err != nil {
			log.Fatal(err)
			return
		}

		defer listen.Close()

	
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("-----### Connection Established ###-----")
		packetChannel := getLinesChannel(conn)
		for line := range packetChannel {
			fmt.Printf("%s\n", line)
		}
	}


	
	
	return
}




