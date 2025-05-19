package request

import (
	"io"
	"log"
	"fmt"
	"unicode"
	"strings"
)

const StateInitialized = 0
const StateDone = 1
const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	State		int
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method			string
}



func RequestFromReader (reader io.Reader) (*Request, error) {
	
	// Initailzie the new buffer for the reader
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	// Create a NEW request with inital state
	r := &Request{
		State: StateInitialized,
		}

	// A "for" (really a while) loop to keep going until done
	for r.State != StateDone {

		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				// We've reached the end of the input
				r.State = StateDone
				break
			}
			return nil, err
		}

		readToIndex += n

		// Use parse method to process the data
		bytesProcessed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Remove processed bytes from buffer
		copy(buf, buf[bytesProcessed:readToIndex])
		readToIndex -= bytesProcessed

	}
	return r, nil
}

func parseRequestLine (input string) (RequestLine, int, error) {

	full := strings.Index(input, "\r\n")
	if full == -1 {
		return RequestLine{}, 0, nil
	}

	lines := strings.Split(input, "\r\n")
	if len(lines) == 0 {
		return RequestLine{}, 0, fmt.Errorf("empty request")
	}


	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return RequestLine{}, 0, fmt.Errorf("incorrect amount of request line")
	}

	request := RequestLine{
		HttpVersion: requestLine[2],
		RequestTarget: requestLine[1],
		Method: requestLine[0],
	}

	if request.HttpVersion != "HTTP/1.1" {
		return RequestLine{}, 0, fmt.Errorf("Incorrect HTTP version or format")
	}

	request.HttpVersion = strings.TrimPrefix(request.HttpVersion, "HTTP/")

	if request.Method == "" {
		return RequestLine{}, 0, fmt.Errorf("empty method")
	}

	if request.RequestTarget == "" {
		return RequestLine{}, 0, fmt.Errorf("empty target, where are you wanting to go?")
	}

	for _, r := range request.Method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return RequestLine{}, 0, fmt.Errorf("Incorrect case or rune for method")
		}
	}	
	return request, full + 2, nil

}

func (r *Request) parse(data []byte) (int, error) {

	if r.State == StateDone {
		return 0, fmt.Errorf("State complete, already finished parsing")
	}

	parsedLine, num, err := parseRequestLine((string(data)))
	if err != nil {
		fmt.Print("unable to parse Request Line")
		return 0, err
	}

	if num == 0 {
		log.Print("Waiting for more data")
		return 0, nil
	}
	
	r.RequestLine = parsedLine
	r.State = StateDone
	return num, nil

}




