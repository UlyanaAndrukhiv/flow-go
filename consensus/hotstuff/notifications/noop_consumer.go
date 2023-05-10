package notifications

import (
	"time"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/model/flow"
)

// NoopConsumer is an implementation of the notifications consumer that
// doesn't do anything.
type NoopConsumer struct {
	NoopProposalViolationConsumer
	NoopFinalizationConsumer
	NoopPartialConsumer
	NoopCommunicatorConsumer
}

var _ hotstuff.Consumer = (*NoopConsumer)(nil)

func NewNoopConsumer() *NoopConsumer {
	nc := &NoopConsumer{}
	return nc
}

// no-op implementation of hotstuff.Consumer(but not nested interfaces)

type NoopPartialConsumer struct{}

func (*NoopPartialConsumer) OnEventProcessed() {}

func (*NoopPartialConsumer) OnStart(uint64) {}

func (*NoopPartialConsumer) OnReceiveProposal(uint64, *model.Proposal) {}

func (*NoopPartialConsumer) OnReceiveQc(uint64, *flow.QuorumCertificate) {}

func (*NoopPartialConsumer) OnReceiveTc(uint64, *flow.TimeoutCertificate) {}

func (*NoopPartialConsumer) OnPartialTc(uint64, *hotstuff.PartialTcCreated) {}

func (*NoopPartialConsumer) OnLocalTimeout(uint64) {}

func (*NoopPartialConsumer) OnViewChange(uint64, uint64) {}

func (*NoopPartialConsumer) OnQcTriggeredViewChange(uint64, uint64, *flow.QuorumCertificate) {}

func (*NoopPartialConsumer) OnTcTriggeredViewChange(uint64, uint64, *flow.TimeoutCertificate) {}

func (*NoopPartialConsumer) OnStartingTimeout(model.TimerInfo) {}

func (*NoopPartialConsumer) OnCurrentViewDetails(uint64, uint64, flow.Identifier) {}

// no-op implementation of hotstuff.FinalizationConsumer

type NoopFinalizationConsumer struct{}

var _ hotstuff.FinalizationConsumer = (*NoopFinalizationConsumer)(nil)

func (*NoopFinalizationConsumer) OnBlockIncorporated(*model.Block) {}

func (*NoopFinalizationConsumer) OnFinalizedBlock(*model.Block) {}

// no-op implementation of hotstuff.TimeoutCollectorConsumer

type NoopTimeoutCollectorConsumer struct{}

var _ hotstuff.TimeoutCollectorConsumer = (*NoopTimeoutCollectorConsumer)(nil)

func (*NoopTimeoutCollectorConsumer) OnTcConstructedFromTimeouts(*flow.TimeoutCertificate) {}

func (*NoopTimeoutCollectorConsumer) OnPartialTcCreated(uint64, *flow.QuorumCertificate, *flow.TimeoutCertificate) {
}

func (*NoopTimeoutCollectorConsumer) OnNewQcDiscovered(*flow.QuorumCertificate) {}

func (*NoopTimeoutCollectorConsumer) OnNewTcDiscovered(*flow.TimeoutCertificate) {}

func (*NoopTimeoutCollectorConsumer) OnTimeoutProcessed(*model.TimeoutObject) {}

// no-op implementation of hotstuff.CommunicatorConsumer

type NoopCommunicatorConsumer struct{}

var _ hotstuff.CommunicatorConsumer = (*NoopCommunicatorConsumer)(nil)

func (*NoopCommunicatorConsumer) OnOwnVote(flow.Identifier, uint64, []byte, flow.Identifier) {}

func (*NoopCommunicatorConsumer) OnOwnTimeout(*model.TimeoutObject) {}

func (*NoopCommunicatorConsumer) OnOwnProposal(*flow.Header, time.Time) {}

// no-op implementation of hotstuff.VoteCollectorConsumer

type NoopVoteCollectorConsumer struct{}

var _ hotstuff.VoteCollectorConsumer = (*NoopVoteCollectorConsumer)(nil)

func (*NoopVoteCollectorConsumer) OnQcConstructedFromVotes(*flow.QuorumCertificate) {}

func (*NoopVoteCollectorConsumer) OnVoteProcessed(*model.Vote) {}

// no-op implementation of hotstuff.ProposalViolationConsumer

type NoopProposalViolationConsumer struct{}

var _ hotstuff.ProposalViolationConsumer = (*NoopProposalViolationConsumer)(nil)

func (*NoopProposalViolationConsumer) OnInvalidBlockDetected(model.InvalidBlockError) {}

func (*NoopProposalViolationConsumer) OnDoubleProposeDetected(*model.Block, *model.Block) {}

func (*NoopProposalViolationConsumer) OnDoubleVotingDetected(*model.Vote, *model.Vote) {}

func (*NoopProposalViolationConsumer) OnInvalidVoteDetected(model.InvalidVoteError) {}

func (*NoopProposalViolationConsumer) OnVoteForInvalidBlockDetected(*model.Vote, *model.Proposal) {}

func (*NoopProposalViolationConsumer) OnDoubleTimeoutDetected(*model.TimeoutObject, *model.TimeoutObject) {
}

func (*NoopProposalViolationConsumer) OnInvalidTimeoutDetected(model.InvalidTimeoutError) {}
