package main

import (
	"net"
	"log"
	"bufio"
	"os"
	"fmt"
)

const network = "udp"
const address = "localhost:42069"


func main() {

	// Start new UDP address connection
	resolver, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		log.Fatal(err)
		return
	}

	// New UDP dial?

	conn, err := net.DialUDP(network, nil, resolver)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		byteLine := []byte(line)
		_, err = conn.Write(byteLine)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
			

