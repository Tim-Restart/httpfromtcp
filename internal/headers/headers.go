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
	h := make(Headers)
	return h
}


func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	
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

		keyWhiteSpace := fullHeader[:colonIndex] 
		valueWhiteSpace := fullHeader[colonIndex +1:]

		// Didn't include backslash.... not sure how too, tried single ''
		if bytes.ContainsAny(keyWhiteSpace, "(),/:;<=>?@[]{})") == true {
			return 0, false, fmt.Errorf("Invalid character foudn in header")
		}
		keySpace := string(keyWhiteSpace)
		valueSpace := string(valueWhiteSpace)
		key := strings.TrimSpace(strings.ToLower(keySpace))
		value := strings.TrimSpace(valueSpace)

		// Checks if the key exists then adds
		_, exists := h[key]
		if exists {
			fmt.Println("Entered the exists condition")
			h[key] += ", " + value
			return n + 2, false, nil
		}


		h[key] = value
		return n + 2, false, nil
	}

}


 


		