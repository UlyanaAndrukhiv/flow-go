package unicast

import (
	"golang.org/x/time/rate"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/onflow/flow-go/network/message"
)

// StreamsRateLimiterImpl unicast rate limiter that limits the amount of streams that can
// be created per some configured interval. A new stream is created each time a libP2P
// node sends a direct message.
type StreamsRateLimiterImpl struct {
	rateLimitedPeers *rateLimitedPeersMap
	limiters         *rateLimiterMap
	limit            rate.Limit
	burst            int
	now              GetTimeNow
}

// NewStreamsRateLimiter returns a new StreamsRateLimiterImpl. The cleanup loop will be started in a
// separate goroutine and should be stopped by calling Close.
func NewStreamsRateLimiter(limit rate.Limit, burst int, now GetTimeNow) *StreamsRateLimiterImpl {
	l := &StreamsRateLimiterImpl{
		rateLimitedPeers: newRateLimitedPeersMap(rateLimiterTTL, cleanUpTickDuration),
		limiters:         newLimiterMap(rateLimiterTTL, cleanUpTickDuration),
		limit:            limit,
		burst:            burst,
		now:              now,
	}

	go l.limiters.cleanupLoop()
	go l.rateLimitedPeers.cleanupLoop()

	return l
}

// Allow checks the cached limiter for the peer and returns limiter.Allow().
// If a limiter is not cached for a one is created.
func (s *StreamsRateLimiterImpl) Allow(peerID peer.ID, _ *message.Message) bool {
	limiter := s.getLimiter(peerID)
	if !limiter.AllowN(s.now(), 1) {
		s.rateLimitedPeers.store(peerID)
		return false
	} else {
		s.rateLimitedPeers.remove(peerID)
		return true
	}
}

// IsRateLimited returns true is a peer is currently rate limited.
func (s *StreamsRateLimiterImpl) IsRateLimited(peerID peer.ID) bool {
	return s.rateLimitedPeers.exists(peerID)
}

// Close sends cleanup signal to underlying rate limiters and rate limited peers maps. After the rate limiter
// is closed it can not be reused.
func (s *StreamsRateLimiterImpl) Close() {
	s.limiters.close()
	s.rateLimitedPeers.close()
}

// getLimiter returns limiter for the peerID, if a limiter does not exist one is created and stored.
func (s *StreamsRateLimiterImpl) getLimiter(peerID peer.ID) *rate.Limiter {
	if limiter, ok := s.limiters.get(peerID); ok {
		return limiter
	}

	limiter := rate.NewLimiter(s.limit, s.burst)
	s.limiters.store(peerID, limiter)

	return limiter
}
