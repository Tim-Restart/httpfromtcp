package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {

	message, err := os.Open(inputFilePath)
	if err != nil {
		log.Print("error reading %s\n", inputFilePath, err)
		panic(err)
	}
	// for loop tracking 8 bytes
	// print the 8 bytes
	// keep going till non left
	//fmt.Printf("%v", message)
	buff := make([]byte, 8)

	
	
	var currentLine string

	for {
		n, err := message.Read(buff)
		if err == io.EOF {
			// More to read
			fmt.Printf("read: %s\n", currentLine)

			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		word := string(buff[:n])

		splitWord := strings.Split(word, "\n")
		if len(splitWord) > 1 {
			firstWord := splitWord[0]
			secondWord := splitWord[1]
			currentLine += firstWord
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
			currentLine +=  secondWord
			
		} else {
			currentLine += word
			
		}
	
	}
	
	return
}




