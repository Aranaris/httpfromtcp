package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() (headers Headers) {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	d := false

	clrf_regex := regexp.MustCompile(`\r\n`)
	str := string(data)

	loc := clrf_regex.FindStringIndex(str)
	//No CLRF found indicates we need to continue parsing data
	if loc == nil {
		return 0, false, nil
	}
	//CLRF at the start of data, indicates end of headers
	if loc[0] == 0 {
		return 0, true, nil
	}


	trailing_r := regexp.MustCompile(`\s+$`)
	leading_r := regexp.MustCompile(`^\s+`)

	fn, fv, found := strings.Cut(str[0:loc[0]], ":")
	if !found {
		return 0, false, errors.New("missing semicolon, invalid header format")
	}
	if trailing_r.MatchString(fn) {
		return 0, false, errors.New("invalid fieldname format in header")
	}

	fn = leading_r.ReplaceAllString(fn, "")
	fv = leading_r.ReplaceAllString(fv, "")
	fv = trailing_r.ReplaceAllString(fv, "")

	h[fn] = fv

	return len([]byte(str[0:loc[1]])), d, nil
}
