package compressor_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/network/compressor"
)

// TestZstdRoundTrip evaluates that
// 1) reading what has been written by compressor yields in same result
// 2) data is compressed when written.
func TestZstdRoundTrip(t *testing.T) {
	text := "hello world, hello world!"
	textBytes := []byte(text)
	textBytesLen := len(textBytes)
	buf := new(bytes.Buffer)

	zstdComp := compressor.ZstdStreamCompressor{}

	w, err := zstdComp.NewWriter(buf)
	require.NoError(t, err)

	// testing write
	n, err := w.Write(textBytes)
	require.NoError(t, err)
	// written bytes should match original data
	require.Equal(t, n, textBytesLen)
	// written data on buffer should be compressed in size.
	require.Less(t, buf.Len(), textBytesLen)
	require.NoError(t, w.Close())

	// testing read
	r, err := zstdComp.NewReader(buf)
	require.NoError(t, err)

	b := make([]byte, textBytesLen)
	n, err = r.Read(b)
	// we read the entire buffer on reader, so it should not return an error
	require.NoError(t, err)
	// we should read same number of bytes as we've written
	require.Equal(t, n, textBytesLen)
	// we should read what we have written
	require.Equal(t, b, textBytes)
}
