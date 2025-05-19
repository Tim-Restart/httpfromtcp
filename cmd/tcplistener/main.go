package main

import (
	"fmt"
	"log"
	"net"
	"github.com/Tim-Restart/httpfromtcp/internal/request"
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
		newLine, err := RequestFromReader(conn)
		if err != nil {
			fmt.Print("Error reading from request")
		}
		fmt.Printf("Request line:/n- Method: %v/n- Target: %v/n- Version: %v/n", newLine.RequestLine.Method, newLine.RequestLine.RequestTarget, newLine.RequestLine.HttpVersion)
	}


	
	
	return
}




