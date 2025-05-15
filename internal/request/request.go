package request

import (
	"io"
	"log"
	"fmt"
	"unicode"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method			string
}

func RequestFromReader (reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader) 
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rl, err := parseRequestLine(string(r))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Request{RequestLine: rl}, nil
}

func parseRequestLine (input string) (RequestLine, error) {

	lines := strings.Split(input, "\r\n")
	if len(lines) == 0 {
		return RequestLine{}, fmt.Errorf("empty request")
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return RequestLine{}, fmt.Errorf("incorrect amount of request line")
	}

	request := RequestLine{
		HttpVersion: requestLine[2],
		RequestTarget: requestLine[1],
		Method: requestLine[0],
	}

	if request.HttpVersion != "HTTP/1.1" {
		return RequestLine{}, fmt.Errorf("Incorrect HTTP version or format")
	}

	request.HttpVersion = strings.TrimPrefix(request.HttpVersion, "HTTP/")

	if request.Method == "" {
		return RequestLine{}, fmt.Errorf("empty method")
	}

	if request.RequestTarget == "" {
		return RequestLine{}, fmt.Errorf("empty target, where are you wanting to go?")
	}

	for _, r := range request.Method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return RequestLine{}, fmt.Errorf("Incorrect case or rune for method")
		}
	}
	


	return request, nil

}




