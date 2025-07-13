package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Status int
const (
	Initialized Status = iota
	Done
)

func (s Status) String() string {
	switch s {
	case 0:
		return "Initialized"
	case 1:
		return "Done"
	}
	return fmt.Sprintf("Status(%q)", int(s))
}

var BufferSize = 8

type Request struct {
	RequestLine RequestLine
	Status 			Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	var r Request
	buffer := make([]byte, BufferSize)
	r.Status = 0
	readToIndex := 0
	parsed := 0

	for {
		if r.Status == 1 {
			return &r, nil
		}

		if readToIndex >= len(buffer) {
			newSize := len(buffer) + BufferSize
			temp := buffer
			buffer = make([]byte, newSize)
			copy(buffer, temp)
		}

		b, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.Status = 1
			}
			return nil, err
		}
		readToIndex += b
		
		n, err := r.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		if n > 0 {
			parsed += n
			newSlice := make([]byte, BufferSize)
			copy(newSlice, buffer[parsed:])
			buffer = newSlice
			readToIndex -= parsed
		}

	}
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	var rq RequestLine
	str := string(data)
	segments := strings.Split(str, "\r\n")
	if len(segments) == 1 {
		return rq, 0, nil
	}

	rs := strings.Split(segments[0], " ")
	

	if len(rs) != 3 {
		return rq, 0, fmt.Errorf("invalid number of parts in request line")
	}

	var f bool
	rq.HttpVersion, f = strings.CutPrefix(rs[2], "HTTP/")
	if !f {
		return rq, 0, fmt.Errorf("invalid HttpVersion syntax")
	}

	r, err := regexp.Compile("([A-Z]+)")
	if err != nil {
		return rq, 0, err
	}
	if !r.MatchString(rs[0]) {
		return rq, 0, fmt.Errorf("invalid request method")
	}

	rq.Method = rs[0]
	rq.RequestTarget = rs[1]
	return rq, len([]byte(segments[0])), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status != 0 {
		return 0, fmt.Errorf("unable to parse request, not in initialized state")
	}

	rq, b, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if b == 0 {
		return 0, nil
	}

	r.RequestLine = rq
	r.Status = 1
	return b, nil
}
