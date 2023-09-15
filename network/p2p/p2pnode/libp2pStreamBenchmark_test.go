package p2pnode_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/irrecoverable"
	mockmodule "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/network/internal/p2pfixtures"
	"github.com/onflow/flow-go/network/p2p"
	p2ptest "github.com/onflow/flow-go/network/p2p/test"
	"github.com/onflow/flow-go/utils/unittest"
)

func BenchmarkStreamCompression(b *testing.B) {
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
		name := fmt.Sprintf("BenchmarkStreamCompression_%d_%d", size.start, size.end)
		start, end := size.start, size.end

		b.Run(name, func(b *testing.B) {
			unicastOverStream(b, start, end)
			//unicastOverStream(b, start, end, p2ptest.WithPreferredUnicasts([]protocols.ProtocolName{protocols.GzipCompressionUnicast}))
			//unicastOverStream(b, start, end, p2ptest.WithPreferredUnicasts([]protocols.ProtocolName{protocols.ZstdCompressionUnicast}))
		})
	}
}

func unicastOverStream(tb *testing.B, minSize int, maxSize int, opts ...p2ptest.NodeFixtureParameterOption) {
	ctx, cancel := context.WithCancel(context.Background())
	signalerCtx := irrecoverable.NewMockSignalerContext(tb, ctx)

	// Creates nodes
	sporkId := unittest.IdentifierFixture()
	idProvider := mockmodule.NewIdentityProvider(tb)
	streamHandler1, inbound1 := p2ptest.StreamHandlerFixture(tb)
	node1, id1 := p2ptest.NodeFixture(

		tb,
		sporkId,

		tb.Name(),
		idProvider,
		append(opts, p2ptest.WithDefaultStreamHandler(streamHandler1))...)

	streamHandler2, inbound2 := p2ptest.StreamHandlerFixture(tb)
	node2, id2 := p2ptest.NodeFixture(

		tb,
		sporkId,

		tb.Name(),
		idProvider,
		append(opts, p2ptest.WithDefaultStreamHandler(streamHandler2))...)
	ids := flow.IdentityList{&id1, &id2}
	nodes := []p2p.LibP2PNode{node1, node2}
	for i, node := range nodes {
		idProvider.On("ByPeerID", node.Host().ID()).Return(ids[i], true).Maybe()

	}
	p2ptest.StartNodes(tb, signalerCtx, nodes, 10*time.Millisecond)
	defer p2ptest.StopNodes(tb, nodes, cancel, 10*time.Millisecond)

	p2ptest.LetNodesDiscoverEachOther(tb, ctx, nodes, ids)

	for i := 0; i < tb.N; i++ {
		msg, _ := longStringMessageFactoryFixture(minSize, maxSize)
		p2pfixtures.EnsureMessageExchangeOverUnicast(
			tb,
			ctx,
			nodes,
			[]chan string{inbound1, inbound2},
			msg)
	}
}

func longStringMessageFactoryFixture(minSize, maxSize int) (func() string, int) {
	var msgSize int
	return func() string {
		msgSize = rand.Intn(maxSize-minSize+1) + minSize
		randomMsg := generateRandomMessage(msgSize)

		return fmt.Sprintf("%s %d \n", randomMsg, time.Now().UnixNano())
	}, msgSize
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
