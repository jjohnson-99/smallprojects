package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/jjohnson-99/httpFromTCP/internal/request"
)

const port = ":8080"

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Listening for TCP traffic on port %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from ", conn.RemoteAddr())
		
		req, err := request.RequestFromReader(conn)
		if err != nil{
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf("Request line:\n - Method: %s\n - Target: %s\n - Version: %s\n", req.RequestLine.Method,
																				    req.RequestLine.RequestTarget,
																				    req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Printf(" - %s: %s\n", key, val)
		}

		fmt.Println("Body:")
		fmt.Printf("%s", string(req.Body))
	}
}

