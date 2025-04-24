package main

import (
	"fmt"
	"log"
	"os"
	"io"
)

func main() {

	message, err := os.Open("messages.txt")
	if err != nil {
		log.Print("error reading messages.txt")
		panic(err)
	}
	// for loop tracking 8 bytes
	// print the 8 bytes
	// keep going till non left
	//fmt.Printf("%v", message)
	buff := make([]byte, 8)

	
	
	// figure out the total length of the messages in bytes

	for {
		n, err := message.Read(buff)
		if err == io.EOF {
			// More to read
			
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
	if n > 0 {
		word := string(buff[:n])
		fmt.Println("read:", word)
	}
	}
	
	return
}

