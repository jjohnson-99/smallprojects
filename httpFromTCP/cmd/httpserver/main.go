package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"github.com/jjohnson-99/httpFromTCP/internal/request"
	"github.com/jjohnson-99/httpFromTCP/internal/response"
	"github.com/jjohnson-99/httpFromTCP/internal/server"
)

const port = 8080

func handler(w *response.Writer, req *request.Request) { //*server.HandlerError {
	//var handlerError server.HandlerError
	htmlBody := "<html>\n" + 
				"\t<head>\n" + 
				"\t\t<title>%d %s</title>\n" +
				"\t</head>\n" +
				"\t<body>\n" +
				"\t\t<h1>%s</h1>\n" +
				"\t\t<p>%s</p>\n" +
				"\t</body>\n" +
				"</html>"
	switch (req.RequestLine.RequestTarget) {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusCodeBadRequest)
		responseText := fmt.Sprintf(htmlBody, response.StatusCodeBadRequest, 
											  "Bad Request",
										  	  "Bad Request",
									  		  "Your request is bad.")
		req.Headers.Set("Content-Length", strconv.Itoa(len(responseText)))
		w.WriteHeaders(req.Headers)
		w.WriteBody([]byte(responseText))
	case "/myproblem":
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		responseText := fmt.Sprintf(htmlBody, response.StatusCodeInternalServerError, 
											  "Internal Server Error",
										  	  "Internal Server Error",
									  		  "Okay, that was my mistake.")
		req.Headers.Set("Content-Length", strconv.Itoa(len(responseText)))
		w.WriteHeaders(req.Headers)
		w.WriteBody([]byte("test body server error"))
	default:
		w.WriteStatusLine(response.StatusCodeSuccess)
		responseText := fmt.Sprintf(htmlBody, response.StatusCodeSuccess,
											  "OK",
										  	  "Success!",
									  		  "Good request..")
		req.Headers.Set("Content-Length", strconv.Itoa(len(responseText))) // <- issue here, need correct string length
		w.WriteHeaders(req.Headers)
		w.WriteBody([]byte(responseText))
	}
	return
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
