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

	rl, err := parseRequestLine(data)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	var r Request
	r.RequestLine = rl
	return &r, nil
}

func parseRequestLine(data []byte) (RequestLine, error) {
	str := string(data)
	segments := strings.Split(str, "\r\n")
	rs := strings.Split(segments[0], " ")
	var rq RequestLine

	if len(rs) != 3 {
		return rq, fmt.Errorf("invalid number of parts in request line")
	}

	var f bool
	rq.HttpVersion, f = strings.CutPrefix(rs[2], "HTTP/")
	if !f {
		return rq, fmt.Errorf("invalid HttpVersion syntax")
	}

	r, err := regexp.Compile("([A-Z]+)")
	if err != nil {
		return rq, err
	}
	if !r.MatchString(rs[0]) {
		return rq, fmt.Errorf("invalid request method")
	}

	rq.Method = rs[0]
	rq.RequestTarget = rs[1]
	return rq, nil
}
