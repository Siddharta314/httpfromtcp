package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	// Headers     map[string]string
	// Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState int

const (
	stateInitialized parserState = iota //0
	stateDone                           //1
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{state: stateInitialized}
	buf := make([]byte, bufferSize)
	var leftover []byte
	// readCount := 0
	for r.state != stateDone {
		// n, err := reader.Read(buf[readCount:])
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				r.state = stateDone
				break
			}
			return nil, err
		}
		// readCount += n
		dataToParse := append(leftover, buf[:n]...)
		consumed, err := r.parse(dataToParse)
		// consumed, err := r.parse(buf[:readCount])
		if err != nil {
			return nil, err
		}
		// if consumed > 0 {
		// 	copy(buf, buf[consumed:readCount])
		// 	readCount -= consumed
		// }

		// if readCount == len(buf) && consumed == 0 {
		// 	newBuf := make([]byte, len(buf)*2)
		// 	copy(newBuf, buf)
		// 	buf = newBuf
		// }
		leftover = dataToParse[consumed:]
	}
	fmt.Printf("Parsed Request: %s %s %s\n",
		r.RequestLine.Method,
		r.RequestLine.RequestTarget,
		r.RequestLine.HttpVersion,
	)
	return r, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}
	requestLine, err := newRequestLine(string(data[:idx]))
	if err != nil {
		return 0, nil, err
	}
	return idx + len(crlf), requestLine, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		n, request, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *request
		r.state = stateDone
		return n, nil
	case stateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state: %d", r.state)
	}
}

func newRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Expect: 3 parts, got: %d", len(parts))
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	target := parts[1]
	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("Expect: HTTP/1.1, got: %s", parts[2])
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("invalid version: %s", versionParts)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}

/*
func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	dataStr := string(data)
	idx := strings.Index(dataStr, crlf)
	if idx == -1 {
		return nil, fmt.Errorf("CRLF not found")
	}
	return newRequestLine(dataStr[:idx])
}

func newRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Expect: 3 parts, got: %d", len(parts))
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	target := parts[1]
	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("Expect: HTTP/1.1, got: %s", parts[2])
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("invalid version: %s", versionParts)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}
*/
