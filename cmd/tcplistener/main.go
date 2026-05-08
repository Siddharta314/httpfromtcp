package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Siddharta314/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err.Error())
		}
		fmt.Println("Connection accepted:  ", conn.RemoteAddr())
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error parsing request: %v", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)
		fmt.Println("Connection closed: ", conn.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	currentLine := ""
	go func() {
		defer f.Close()
		defer close(ch)

		buffer := make([]byte, 8)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			newLine := string(buffer[:n])
			parts := strings.Split(newLine, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return ch
}
