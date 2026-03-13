package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}


func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	text := string(data)
	requestLine, err := parseRequestLine(text)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(text string) (*RequestLine, error) {
	SEP := "\r\n"
	idx := strings.Index(text, SEP)
	if idx == -1 {
		return nil, fmt.Errorf("invalid request line")
	}
	parts := strings.Split(text[:idx], " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid http version")
	}
	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}, nil
}