package response

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/jjohnson-99/httpFromTCP/internal/headers"
)

type StatusCode int

const (
	StatusCodeSuccess 			  StatusCode = 200
	StatusCodeBadRequest 		  StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type Writer struct {
	Buffer *bytes.Buffer	
}

func (w *Writer) WriteStatusLine(StatusCode StatusCode) error {
	var p []byte
	switch (StatusCode) {
	case StatusCodeSuccess:
		p = []byte("HTTP/1.1 200 OK\r\n")
	case StatusCodeBadRequest:	
		p = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusCodeInternalServerError:
		p = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		p = []byte(fmt.Sprintf("HTTP/1.1 %d\r\n", StatusCode))
	}

	_, err := w.Buffer.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var p []byte

	for name, value := range headers {
		p = []byte(fmt.Sprintf("%s: %s\r\n", name, value))
		_, err := w.Buffer.Write(p)
		if err != nil {
			return err
		}
	}

	w.Buffer.Write([]byte("\r\n"))
	return nil
}


func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Buffer.Write(p)
	if err != nil {
		return 0, nil
	}

	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	body := fmt.Sprintf("%x\r\n", len(p)) + string(p) + "\r\n"
	n, err := w.Buffer.Write([]byte(body))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return 0, nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(map[string]string) //headers.Headers

	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"
	
	return headers
}


