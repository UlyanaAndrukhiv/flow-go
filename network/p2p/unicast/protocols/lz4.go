package protocols

import (
	libp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/network/compressor"
	"github.com/onflow/flow-go/network/p2p/compressed"
)

const Lz4CompressionUnicast = ProtocolName("lz4-compression")

func FlowLz4ProtocolId(sporkId flow.Identifier) protocol.ID {
	return protocol.ID(FlowLibP2PProtocolLz4CompressedOneToOne + sporkId.String())
}

// Lz4Stream is a stream compression creates and returns a lz4-compressed stream out of input stream.
type Lz4Stream struct {
	protocolId     protocol.ID
	defaultHandler libp2pnet.StreamHandler
	logger         zerolog.Logger
}

func NewLz4CompressedUnicast(logger zerolog.Logger, sporkId flow.Identifier, defaultHandler libp2pnet.StreamHandler) *Lz4Stream {
	return &Lz4Stream{
		protocolId:     FlowLz4ProtocolId(sporkId),
		defaultHandler: defaultHandler,
		logger:         logger.With().Str("subsystem", "lz4-unicast").Logger(),
	}
}

// UpgradeRawStream wraps lz4 compression and decompression around the plain libp2p stream.
func (l Lz4Stream) UpgradeRawStream(s libp2pnet.Stream) (libp2pnet.Stream, error) {
	return compressed.NewCompressedStream(s, compressor.Lz4Compressor{})
}

func (l Lz4Stream) Handler(s libp2pnet.Stream) {
	// converts native libp2p stream to lz4-compressed stream
	s, err := l.UpgradeRawStream(s)
	if err != nil {
		l.logger.Error().Err(err).Msg("could not create compressed stream")
		return
	}
	l.defaultHandler(s)
}

func (l Lz4Stream) ProtocolId() protocol.ID {
	return l.protocolId
}
