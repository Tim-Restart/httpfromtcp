package headers

import (
	"strings"
	"fmt"
	"bytes"
	"log"
)

const CRLF = "\r\n"
const colon = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}


func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	totalConsumed := 0

	for {
		idx := bytes.Index(data, []byte(CRLF))

		if idx == -1 {
			return totalConsumed, false, nil
		}
		if idx == 0 {
			// the empy line
			// headers are done, consume CRLF
			totalConsumed += 2
			return totalConsumed, true, nil
		}

		headerLine := data[:idx]

		parts := bytes.SplitN(headerLine, []byte(":"), 2) // this is the key logic I didn't have
		key := strings.ToLower(string(parts[0])) 

		if key != strings.TrimRight(key, " ") {
			return 0, false, fmt.Errorf("Invalid header name: %s", key)
		}

		value := bytes.TrimSpace(parts[1])
		key = strings.TrimSpace(key)
		if !validTokens([]byte(key)) {
			return 0, false, fmt.Errorf("invalid header token found: %s", key)
		}
		h.Set(key, string(value))
		return idx + 2, false, nil
	}
	return totalConsumed, true, nil
}


var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '.', '^', '_', '`', '|', '~'}

func validTokens(data []byte) bool {
	for _, c := range data {
		if !(c >= 'A' && c <= 'Z' ||
			c >= 'a' && c <= 'z' || 
			c >= '0' && c <= '9' ||
			c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func (h Headers) Get(content string) (string, bool) {
	//log.Print("** You are in the Get Headers place **")
	key := strings.ToLower(content)
	value, ok := h[key]
	if ok {
		log.Printf("value: %v", value)
		return value, true
	} else {
		return "", false
	}
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}


 


		