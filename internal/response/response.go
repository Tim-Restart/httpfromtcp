package response

import (
	"fmt"
	"io"
	"httpfromtcp/internal/headers"
	"strconv"
)

type StatusCode int

type Writer struct {
	// I want to put status code in here??
	StatusLine 	[]byte
	Headers		[]byte
	Body 		[]byte
}

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
// Switches on the response code, and sends the HTML response to the w.Writer - need to add writer still
func HtmlResponse(resposneCode StatusCode) {
	switch resposneCode{
	case BadRequest:
		w.Write([]byte("<html>\n  <head>\n\t<title>400 Bad Request</title>\n  </head>\n  <body>\n\t<h1>Bad Request</h1>\n\t<p>Your request honestly kinda sucked.</p>\n  </body>\n</html>"))
		return nil
	case InternalError:
		w.Write([]byte("<html>\n  <head>\n\t<title>500 Internal Server Error</title>\n  </head>\n  <body>\n\t<h1>Internal Server Error</h1>\n\t<p>Okay, you know what? This one is on me.</p>\n  </body>\n</html>"))
		return nil
	default:
		w.Write([]byte("<html>\n  <head>\n\t<title>200 OK</title>\n  </head>\n  <body>\n\t<h1>Success!</h1>\n\t<p>Your request was an absolute banger.</p>\n  </body>\n</html>"))
		return nil
	}
}

