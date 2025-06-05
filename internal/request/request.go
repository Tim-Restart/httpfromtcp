package request

import (
	"io"
	"fmt"
	"unicode"
	"strings"
	"strconv"
	"httpfromtcp/internal/headers"
	"log"
)

const StateInitialized = 0
const StateDone = 1
const RequestStateParsingHeaders = 2
const RequestStateParsingBody = 3
const bufferSize = 8
const content = "content-length"

type Request struct {
	RequestLine RequestLine
	Headers 	headers.Headers // Using the headers type from the headers package
	Body 		[]byte
	bodyLengthRead int
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
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
		}

	// A "for" (really a while) loop to keep going until done
	for r.State != StateDone {

		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		log.Println("entered the part before the EOF part!!!!")
		if err != nil {
			log.Println("entered pre EOF")
			if err == io.EOF {
				log.Println("Entered the EOF")
				if r.State != StateDone {
					return nil, err
				}
				// We've reached the end of the input
				// r.State = StateDone
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
	log.Printf("This is the request: %v", request)	
	return request, full + 2, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {
	remaining := data
	switch r.State {
	case StateInitialized:
		// logic to parse request line
		reqLine, n, err := parseRequestLine(string(remaining))
		
		if err != nil {
			return 0, fmt.Errorf("Error passing Request Line")
		}
		if n == 0 {
			// Just need more data
			return 0, nil
		}
		//remaining = remaining[n:]
		r.RequestLine = reqLine
		r.State = RequestStateParsingHeaders
		log.Printf("This is the Request Line: %v", reqLine)
		return n, nil

	case RequestStateParsingHeaders:
		// logic to parse headers

		n, done, err := r.Headers.Parse(remaining)
		
		if err != nil {
			fmt.Printf("Error: %v", err)
			fmt.Printf("Returning header error: %v", string(data))
			return 0, fmt.Errorf("Error parsing headers")
		}
		if done == true {
			log.Printf("^^^^ This is the headers: %v", r.Headers)
			r.State = RequestStateParsingBody
		}
		remaining = remaining[n:]
		return n, nil

	case RequestStateParsingBody:
		log.Println("----- You have entered the body parsing zone -----")
		number, ok := r.Headers.Get(content)
		if !ok {
			r.State = StateDone
			return len(data), nil
		}
		contentLength, err := strconv.Atoi(number)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		log.Printf("This is the body: %v", r.Body)
		if r.bodyLengthRead > contentLength {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.bodyLengthRead == contentLength {
			r.State = StateDone
		}
		
		return len(data), nil
		

	default:
		return 0, fmt.Errorf("unknown state")
	}
	
	return 0, nil
}


func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.State != StateDone {
		//log.Printf("ParseSingle Call Log Print: %v",string(data[totalBytesParsed:]))
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			// Need more data to make progress in the current state
			break
		}
		
	}
	return totalBytesParsed, nil
}

//func (r *Request) parseBody(data []byte, length int) (int, error) {
//	if len(data) < length {
//		return 0, fmt.Errorf("Data length less than content length")
//	} 
//
//	body := data[:length]
//	r.Body = body
//	return len(body), nil
//}






