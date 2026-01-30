package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = "localhost:8080"

func main() {
	udp, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("error resolving %s: %s\n", port, err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatalf("error preparing UDP connection: %s\n", err.Error())
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading: %s", err.Error())
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("error writing: %s", err.Error())
		}
	}
}
