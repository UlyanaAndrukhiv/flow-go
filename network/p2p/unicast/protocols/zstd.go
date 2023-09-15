package protocols

import (
	libp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/network/compressor"
	"github.com/onflow/flow-go/network/p2p/compressed"
)

const ZstdCompressionUnicast = ProtocolName("zstd-compression")

func FlowZstdProtocolId(sporkId flow.Identifier) protocol.ID {
	return protocol.ID(FlowLibP2PProtocolZstdCompressedOneToOne + sporkId.String())
}

// ZstdStream is a stream compression creates and returns a zstd-compressed stream out of input stream.
type ZstdStream struct {
	protocolId     protocol.ID
	defaultHandler libp2pnet.StreamHandler
	logger         zerolog.Logger
}

func NewZstdCompressedUnicast(logger zerolog.Logger, sporkId flow.Identifier, defaultHandler libp2pnet.StreamHandler) *ZstdStream {
	return &ZstdStream{
		protocolId:     FlowZstdProtocolId(sporkId),
		defaultHandler: defaultHandler,
		logger:         logger.With().Str("subsystem", "zstd-unicast").Logger(),
	}
}

// UpgradeRawStream wraps zstd compression and decompression around the plain libp2p stream.
func (z ZstdStream) UpgradeRawStream(s libp2pnet.Stream) (libp2pnet.Stream, error) {
	return compressed.NewCompressedStream(s, compressor.ZstdStreamCompressor{})
}

func (z ZstdStream) Handler(s libp2pnet.Stream) {
	// converts native libp2p stream to zstd-compressed stream
	s, err := z.UpgradeRawStream(s)
	if err != nil {
		z.logger.Error().Err(err).Msg("could not create compressed stream")
		return
	}
	z.defaultHandler(s)
}

func (z ZstdStream) ProtocolId() protocol.ID {
	return z.protocolId
}
