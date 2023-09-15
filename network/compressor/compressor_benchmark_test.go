package compressor_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/compressor"
	"github.com/onflow/flow-go/network/p2p/unicast/protocols"
)

func BenchmarkStreamCompression(b *testing.B) {
	// Define a list of data size ranges to benchmark
	sizes := []struct {
		start int
		end   int
	}{
		{50000, 100000},
		{100000, 150000},
		{150000, 200000},
		{200000, 250000},

		{250000, 300000},
		{300000, 350000},
		{350000, 400000},
		{400000, 450000},

		{450000, 500000},
		{500000, 550000},
		{550000, 600000},
		{600000, 650000},
	}

	for _, size := range sizes {
		// Generate a benchmark name based on the data size range
		name := fmt.Sprintf("BenchmarkStreamCompression_%d_%d", size.start, size.end)
		start, end := size.start, size.end

		b.Run(name, func(b *testing.B) {
			// Call the roundTrip function to benchmark compression and decompression
			roundTrip(b, start, end, protocols.GzipCompressionUnicast)
			// Uncomment the following lines to benchmark other compression protocols
			//roundTrip(b, start, end, protocols.ZstdCompressionUnicast)
			//roundTrip(b, start, end, protocols.Lz4CompressionUnicast)
		})
	}
}

func roundTrip(b *testing.B, minSize int, maxSize int, protocolName protocols.ProtocolName) {
	var comp network.Compressor

	// Select the compressor based on the specified protocol
	switch protocolName {
	case protocols.GzipCompressionUnicast:
		comp = compressor.GzipStreamCompressor{}
	case protocols.ZstdCompressionUnicast:
		comp = compressor.ZstdStreamCompressor{}
	case protocols.Lz4CompressionUnicast:
		comp = compressor.Lz4Compressor{}
	}

	for i := 0; i < b.N; i++ {
		// Generate a random string message within the specified size range
		text, _ := longStringMessageFactoryFixture(minSize, maxSize)
		textBytes := []byte(text)
		textBytesLen := len(textBytes)
		buf := new(bytes.Buffer)

		w, err := comp.NewWriter(buf)
		require.NoError(b, err)

		// testing write
		n, err := w.Write(textBytes)
		require.NoError(b, err)
		// written bytes should match original data
		require.Equal(b, n, textBytesLen)
		// written data on buffer should be compressed in size.
		require.Less(b, buf.Len(), textBytesLen)
		require.NoError(b, w.Close())

		// Decompress the data
		r, err := comp.NewReader(buf)
		require.NoError(b, err)

		// Read all decompressed data
		decompressedBytes, err := io.ReadAll(r)
		require.NoError(b, err)

		// Verify decompressed data matches the original
		require.Equal(b, textBytes, decompressedBytes)
	}
}

// TODO: move to another common package
// longStringMessageFactoryFixture generates a random string message within the specified size range.
func longStringMessageFactoryFixture(minSize, maxSize int) (string, int) {
	msgSize := rand.Intn(maxSize-minSize+1) + minSize
	randomMsg := generateRandomMessage(msgSize)
	return randomMsg, msgSize
}

// TODO: move to another common package
// generateRandomMessage generates a random string of the specified length.
func generateRandomMessage(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randomBytes)
}
