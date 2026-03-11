package main

import (
	"fmt"
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
  defer file.Close()

  currentLine := ""
  for {
    buffer := make([]byte, 8)
    n, err := file.Read(buffer)
    if err != nil {
      break
    }
    currentLine += string(buffer[:n])
    parts := strings.Split(currentLine, "\n")
    if len(parts) > 1{
      fmt.Printf("read: %s\n", parts[0])
      currentLine = parts[1]
    }
    
  }
}

