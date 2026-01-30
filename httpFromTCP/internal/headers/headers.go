package headers

import (
	"fmt"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) Get(headerName string) string {
	headerName = strings.ToLower(headerName)
	val, ok := h[headerName]
	if !ok {
		return ""
	}
	return val
}

func (h Headers) Set(headerName, headerValue string) {
	headerName = strings.ToLower(headerName)
	h[headerName] = headerValue
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	numBytesParsed := 0

	if len(string(data)) >= 2 && string(data[:2]) == "\r\n" { // \r\n at the beginning of the data so we've reached the end of the headers
		return 0, true, nil
	}
	fieldLines := strings.Split(string(data), "\r\n")
	if len(fieldLines) == 1 { // we haven't hit a \r\n so we need more data
		return 0, false, nil
	}
	numBytesParsed += 2 // corresponding to the first \r\n found

	header := fieldLines[0]
	numBytesParsed += len(header)
	headerParts := strings.Split(string(header), ": ")
	if len(headerParts) != 2 {
		return 0, false, fmt.Errorf("Malformed hearder: %s", header) //errors.New("Malformed header.")
	}
	headerName   := headerParts[0]
	headerValue := headerParts[1]
	if string(headerName[len(headerName)-1]) == " " {
		return 0, false, errors.New("Cannot have whitespace between field-name and ':'.")
	}
	
	allowedCharacters := regexp.MustCompile("^[0-9a-zA-Zi!#$%&'*+-.^_|~`]+$").MatchString
	if !allowedCharacters(headerName) {
		return 0, false, errors.New("Invalid characters in header.")
	}

	headerName   = strings.ToLower(strings.TrimSpace(headerName))
	headerValue = strings.TrimSpace(headerValue)

	val, ok := h[headerName]
	if ok {
		h[headerName] = val + ", " + headerValue
	} else {
		h[headerName] = headerValue
	}

	return numBytesParsed, false, nil
}
