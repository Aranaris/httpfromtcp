package request

import (
	"fmt"
	"io"
	"regexp"
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
		fmt.Printf("%v", err)
		return nil, err
	}

	rl, _, err := parseRequestLine(data)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	var r Request
	r.RequestLine = rl
	return &r, nil
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
	return rq, len(data), nil
}
