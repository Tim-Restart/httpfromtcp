package main

import (
	"fmt"
	"log"
	"net"
	"httpfromtcp/internal/request"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const host = ":42069"




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
		newLine, err := request.RequestFromReader(conn)
		conn.Close()
		if err != nil {
			log.Print("Error reading from request")
			continue
		}
		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", newLine.RequestLine.Method, newLine.RequestLine.RequestTarget, newLine.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range newLine.Headers {
			fmt.Printf("- %v: %v\n", key, value)
		}
		
	}
	return
}




