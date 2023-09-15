package compressor

import (
	"io"

	"github.com/DataDog/zstd"

	"github.com/onflow/flow-go/network"
)

// ZstdStreamCompressor is a struct representing a Zstd stream compressor.
type ZstdStreamCompressor struct {
}

// NewReader creates a new Zstd reader that decompresses data from the given input reader.
func (zstdStreamComp ZstdStreamCompressor) NewReader(r io.Reader) (io.ReadCloser, error) {
	return zstd.NewReader(r), nil
}

// NewWriter creates a new Zstd writer that compresses data and writes it to the given output writer.
func (zstdStreamComp ZstdStreamCompressor) NewWriter(w io.Writer) (network.WriteCloseFlusher, error) {
	return &zstdWriteCloseFlusher{w: zstd.NewWriter(w)}, nil
}

// zstdWriteCloseFlusher is a private struct representing a Zstd writer with Close and Flush methods.
type zstdWriteCloseFlusher struct {
	w *zstd.Writer
}

// Write writes compressed data to the Zstd writer.
func (zstdW *zstdWriteCloseFlusher) Write(p []byte) (int, error) {
	return zstdW.w.Write(p)
}

// Close closes the underlying Zstd writer.
func (zstdW *zstdWriteCloseFlusher) Close() error {
	return zstdW.w.Close()
}

// Flush flushes any buffered data to the underlying writer.
func (zstdW *zstdWriteCloseFlusher) Flush() error {
	return zstdW.w.Flush()
}
