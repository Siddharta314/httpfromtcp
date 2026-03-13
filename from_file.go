package main

import (
	"fmt"
	"os"
)

const inputFilePath = "messages.txt"


func ReadFromFile(path string) *os.File {
  file, err := os.Open(inputFilePath)
  if err != nil {
    fmt.Println("Error opening file:", err)
    return nil
  }
  defer file.Close()
  return file
}