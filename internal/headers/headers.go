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
// Commented out this logic as it didn't seem to be working
	/*
	switch bytes.Index(data, []byte(CRLF)) {
    case 0:
        return 2, true, nil
	case -1:
        return 0, false, nil
    default:
		n := bytes.Index(data, []byte(CRLF))
		fullHeader := data[:n]
		colonIndex := bytes.Index(fullHeader, []byte(colon))
		if colonIndex > 0 && fullHeader[colonIndex -1] == ' ' {
			return 0, false, fmt.Errorf("Invalid header: space before colon")
		}
		if colonIndex == -1 { // Colon not found invalid header
			fmt.Println("Returning header error")
			return 0, false, fmt.Errorf("Invalid header: No colon found")
		}

		keyWhiteSpace := fullHeader[:colonIndex] 
		valueWhiteSpace := fullHeader[colonIndex +1:]
		//endLineCheck := fullHeader[:colonIndex +1]
		
		keySpace := string(keyWhiteSpace)
		valueSpace := string(valueWhiteSpace)
		key := strings.TrimSpace(strings.ToLower(keySpace))
		value := strings.TrimSpace(valueSpace)
		if !validTokens([]byte(key)) {
			return 0, false, fmt.Errorf("invalid header token found: %s", key)
		}

		// Checks if the key exists then adds
		//_, exists := h[key]
		//if exists {
		//	fmt.Println("Entered the exists condition")
		//	h[key] += ", " + value
		//	return n + 2, false, nil
		//}
	
		h.Set(key, string(value))
		return n + 2, false, nil
	}

}
*/

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
	log.Print("** You are in the Get Headers place **")
	key := strings.ToLower(content)
	value, ok := h[key]
	if ok {
		log.Printf("value: %v", value)
		return value, true
	} else {
		log.Printf("No value in Get")
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


 


		