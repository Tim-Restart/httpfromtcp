package response

import (
	"fmt"
	"io"
	"httpfromtcp/internal/headers"
	"strconv"
)

type StatusCode int



type Writer struct {
	writer			io.Writer
	WriteStatus		int
}

const (
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalError StatusCode = 500
	crlf = "\r\n"
	StatusLine = 0
	WriteHeader = 1
	WriteBody = 2
	Finished = 3
)



func NewWriter( w io.Writer) *Writer{
	newWriter := &Writer{
		writer: w,
		WriteStatus: 0,
	}
	return newWriter
}


func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	// Initializes the Writer status for the Status line
	if w.WriteStatus != 0 {
		return fmt.Errorf("Not ready to receive status line yet")
	}
	
	switch statusCode{
	case Ok:
		_, err := w.writer.Write([]byte("HTTP/1.1 200 OK" + crlf))
		if err != nil {
			fmt.Println(err)
			return err
		}
	
		w.WriteStatus = 1
		return nil
	case BadRequest:
		_, err := w.writer.Write([]byte("HTTP/1.1 400 Bad Request" + crlf))
		if err != nil {
			fmt.Println(err)
			return err
		}
	
		w.WriteStatus = 1
		return nil
	case InternalError:
		_, err := w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error" + crlf))
		if err != nil {
			fmt.Println(err)
			return err
		}
	
		w.WriteStatus = 1
		return nil
	default:
		w.WriteStatus = 1
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


func (w *Writer) WriteHeaders(headers headers.Headers) error {
	// Checks if the writer is in the correct status for headers
	if w.WriteStatus != 1 {
		return fmt.Errorf("Not ready to receive headers yet")
	}


	for k, v := range headers {
		joined := k + ": " + v + "\r\n"
		formatted := []byte(joined)
		_, err := w.writer.Write(formatted)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	endHeader := []byte("\r\n")
	_, err := w.writer.Write(endHeader)
	if err != nil {
			fmt.Println(err)
			return err
	}
	w.WriteStatus = 2
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriteStatus != 2 {
		return 0, fmt.Errorf("Not ready to receive body yet")
	}
	length := len(p)
	_, err := w.writer.Write(p)
	if err != nil {
			fmt.Println(err)
			return 0, err
	}
	// Reset the write status to finihshed
	w.WriteStatus = 3
	return length, nil

}

// Switches on the response code, and sends the HTML response to the w.Writer - need to add writer still
func HtmlResponse(responseCode StatusCode) string{
	switch responseCode{
	case BadRequest:
		badRequest := "<html>\n  <head>\n\t<title>400 Bad Request</title>\n  </head>\n  <body>\n\t<h1>Bad Request</h1>\n\t<p>Your request honestly kinda sucked.</p>\n  </body>\n</html>"
		return badRequest
	case InternalError:
		internalError := "<html>\n  <head>\n\t<title>500 Internal Server Error</title>\n  </head>\n  <body>\n\t<h1>Internal Server Error</h1>\n\t<p>Okay, you know what? This one is on me.</p>\n  </body>\n</html>"
		return internalError
	default:
		okResponse := "<html>\n  <head>\n\t<title>200 OK</title>\n  </head>\n  <body>\n\t<h1>Success!</h1>\n\t<p>Your request was an absolute banger.</p>\n  </body>\n</html>"
		return okResponse
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	// continually parses the chuncked body
	// Chunk sizes should be the sizes in bytes of the data, 
	// and should be in hexadecimal format
	length := len(p)
	hexLength := strconv.FormatInt(int64(length), 16)
	_, err := w.writer.Write([]byte(hexLength + crlf))
	if err != nil {
		return 0, fmt.Errorf("Error writing the length of chuncked body")
	}
	_, err = w.writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("Error writing the body of chuncked body")
	}
	_, err = w.writer.Write([]byte(crlf))
	if err != nil {
		return 0, fmt.Errorf("Error writing the end of chuncked package")
	}
	return length, nil

}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	// does what is required once the chuncked body is done
	endChunk := []byte("0\r\n")
	length, err := w.writer.Write(endChunk)
	if err != nil {
		return 0, fmt.Errorf("Error writing the end of chuncked package")
	}
	return length, nil

}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	// Trailers are formated like headers:
	// Trailer: Lane, Prime, TJ
	// Then at the bottom:
	// Lane: goober
	// Prime: chill-guy
	// TJ: 1-indexer
	// \r\n
	endChunk := []byte(crlf)
	for k, v := range h {
		joined := k + ": " + v + "\r\n"
		formatted := []byte(joined)
		_, err := w.writer.Write(formatted)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	_, err := w.writer.Write(endChunk)
	if err != nil {
			fmt.Println(err)
			return err
		}
	return nil
}


