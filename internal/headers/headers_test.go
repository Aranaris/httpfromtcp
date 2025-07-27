package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid Done
	headers = NewHeaders()
	data = []byte("\r\n   Host:localhost\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("         Host:    localhost:1337      \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:1337", headers["host"])
	assert.Equal(t, 40, n)
	assert.False(t, done)
	
	// Test: Valid 2 headers with existing headers
	data = []byte("Field1: Value1\r\nField2: Value2\r\n\r\n")
	i, done := 0, false
	for !done {
		n, done, err = headers.Parse(data[i:])
		if err != nil {
			break
		}
		i += n
	}
	require.NoError(t, err)
	assert.Equal(t, "localhost:1337", headers["host"])
	assert.Equal(t, "Value1", headers["field1"])
	assert.Equal(t, "Value2", headers["field2"])
	assert.True(t, done)

	// Test: Invalid characters in header key
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
