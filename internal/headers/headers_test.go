package headers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
	"fmt"
)

func TestParse(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	fmt.Printf("Here its the Headers: %v\n", headers)
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

	// Test: Upper case to Lower Case
	headers = NewHeaders()
	data = []byte("HOST: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)


	// Test: for special characters in the key
	headers = NewHeaders()
	data = []byte("HO<ST: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	fmt.Printf("Error recv: %v\n", err)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid pre-existing header
	headers = NewHeaders()
	headers["host"] = "localhost:42069"
	data = []byte("Host: localhost:42067\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:42067", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}