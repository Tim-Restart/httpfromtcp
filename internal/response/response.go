package response

import (
	"fmt"
	"io"
	"httpfromtcp/internal/headers"
	"strconv"
)

type StatusCode int

const (
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalError StatusCode = 500
	crlf = "\r\n"
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode{
	case Ok:
		w.Write([]byte("HTTP/1.1 200 OK" + crlf))
		return nil
	case BadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request" + crlf))
		return nil
	case InternalError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error" + crlf))
		return nil
	default:
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}


func WriteHeaders(w io.Writer, headers headers.Headers) error {
	
	for k, v := range headers {
		joined := k + ": " + v + "\r\n"
		formatted := []byte(joined)
		_, err := w.Write(formatted)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	endHeader := []byte("\r\n")
	_, err := w.Write(endHeader)
	if err != nil {
			fmt.Println(err)
			return err
	}
	return nil
}