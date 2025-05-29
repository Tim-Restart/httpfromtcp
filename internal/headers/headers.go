package headers

import (
	"strings"
	"fmt"
	"bytes"
)

const CRLF = "\r\n"
const colon = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}


func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	
	//if len(data) == 0 {
	//	return 0, false, fmt.Errorf("Error: Empty Header")
	//}

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
	key := strings.ToLower(content)
	value, ok := h[key]
	if ok {
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


 


		