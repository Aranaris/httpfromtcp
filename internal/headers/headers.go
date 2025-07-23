package headers

type Headers map[string]string

func NewHeaders() (headers Headers) {
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	return 0, false, nil
}
