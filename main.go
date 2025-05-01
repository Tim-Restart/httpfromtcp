package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"strings"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)
	
	go func() {
		/* 
		do the reading here? 
		This should be the anon function that does the reading
		*/
		defer f.Close()
		defer close(ch)
		buff := make([]byte, 8)
		var currentLine string
		
		for {
			n, err := f.Read(buff)
			
			if err == io.EOF {
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



func main() {

	message, err := os.Open(inputFilePath)
	if err != nil {
		log.Print("error reading %s\n", inputFilePath, err)
		panic(err)
	}

	packetChannel := getLinesChannel(message)
	for line := range packetChannel {
		fmt.Printf("read: %s\n", line)
	}
	
	return
}




