package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)


func main(){
  listener, err := net.Listen("tcp", ":42069")
  if err != nil {
    fmt.Println("Error starting server:", err)
    return
  }
  defer listener.Close()

  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection:", err)
      return
    }
    defer conn.Close()

    for line := range getLinesChannel(conn) {
      fmt.Println("read:", line)
    }
  }
}


func getLinesChannel(f io.ReadCloser) <-chan string{
  ch := make(chan string)
  currentLine := ""
  go func() {
    defer f.Close()
    defer close(ch)

    buffer := make([]byte, 8)
    for {
      n, err := f.Read(buffer)
      if err != nil {
        break
      }
      currentLine += string(buffer[:n])
      parts := strings.Split(currentLine, "\n")
      if len(parts) > 1{
        ch <- parts[0]
        currentLine = parts[1]
      }
    }
  }()
  return ch
}