package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main(){
  file, err := os.Open(inputFilePath)
  if err != nil {
    fmt.Println("Error opening file:", err)
    return
  }

  lines := getLinesChannel(file)
  for line := range lines {
    fmt.Println("read:", line)
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