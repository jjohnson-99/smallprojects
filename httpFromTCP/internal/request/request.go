package request

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jjohnson-99/httpFromTCP/internal/headers"
)

type requestState int

const (
	requestStateInitialized    requestState = 0
	requestStateParsingHeaders requestState = 1
	requestStateParsingBody    requestState = 2
	requestStateDone           requestState = 3
)

const bufferSize = 8

type Request struct {
	RequestLine   RequestLine
	Headers       headers.Headers
	Body 		  []byte
	state 		  requestState
	contentLength int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method 		  string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		Headers: make(map[string]string),
		Body: make([]byte, 0, bufferSize),
		state: requestStateInitialized,
		contentLength: 0,
	}

	for req.state != requestStateDone {
		if len(buf) == cap(buf) {
			newBuf := make([]byte, 2*len(buf), 2*cap(buf))
			_ = copy(newBuf, buf)
			buf = newBuf
		}

		// read into buffer
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = requestStateDone
				break;
			}
			return nil, err
		}
		readToIndex += numBytesRead

		// parse the buffer
		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		newBuf := make([]byte, len(buf), cap(buf))
		_ = copy(newBuf, buf[numBytesParsed:])
		buf = newBuf
		readToIndex -= numBytesParsed
	}

	if req.contentLength != len(req.Body) {
		return nil, errors.New("Content length mismatch.")
	}

	return req, nil
}


func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 { // need to read more data
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	numBytesParsed := 0
	switch (r.state) {
	case requestStateInitialized:
		n, err := parseRequestLine(r, string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		} else {
			numBytesParsed = n
			r.state = requestStateParsingHeaders
		}
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done == true {
			r.state = requestStateParsingBody
			s := r.Headers.Get("content-length")
			if s != "" {
				r.contentLength, err = strconv.Atoi(s)
				if err != nil {
					return 0, err
				}
			}
			n += 2 // consume \r\n after headers
		}
		numBytesParsed = n
	case requestStateParsingBody:
		n, err := r.parseBody(data)
		if err != nil {
			return 0, nil
		}
		numBytesParsed = n
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state.")
	default:
		return 0, errors.New("error: unknown state.")
	}

	return numBytesParsed, nil
}

func (r *Request) parseBody(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if data[0] == '\r' || data[0] == '\n' {
		return 1, nil
	}

	r.Body = append(r.Body, string(data)...) // append([]byte, byte) is variadic, hence ..., handles resizing
	if len(r.Body) == r.contentLength {
		r.state = requestStateDone
	}
	if len(r.Body) > r.contentLength {
		return 0, errors.New("Request body is longer than Content-Length")
	}

	return len(data), nil
}

// make into method
func parseRequestLine(r *Request, request string) (int, error) {
	numBytesParsed := 0

	lines := strings.Split(request, "\r\n")
	if len(lines) == 1 { // needs more data. should be careful with "\r\n" giving len == 1?
		return 0, nil
	}
	numBytesParsed += 2 // accouting for \r\n

	request_line := lines[0]
	numBytesParsed += len(request_line)

	parts := strings.Split(request_line, " ")
	if len(parts) != 3 {
		return 0, errors.New("Malformed HTTP request: Invalid number of parts in request line.")
	}

	r.RequestLine.Method = parts[0]
	r.RequestLine.RequestTarget = parts[1]
	r.RequestLine.HttpVersion = strings.Split(parts[2], "/")[1]

	isUpperAlpha := regexp.MustCompile(`^[A-Z]+$`).MatchString
	if !isUpperAlpha(r.RequestLine.Method) {
		return 0, errors.New("Expected method to contain only capital alphabetic characters.")
	}
	if r.RequestLine.HttpVersion != "1.1" {
		return 0, errors.New("Expected HTTP version to be 1.1.")
	}

	return numBytesParsed, nil
}

