package test

import (
	"testing"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/network/stub"
	"github.com/onflow/flow-go/utils/unittest"
)

// Wrap several protocol.State modules in one to enable epoch builder
// Use Signer from hotstuff integration tests
func TestClusterSwitchover(t *testing.T) {

	*unittest.LogVerbose = true

	collector := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	participants := unittest.CompleteIdentitySet(collector)
	rootSnapshot := unittest.RootSnapshotFixture(participants)
	log := unittest.Logger()
	hub := stub.NewNetworkHub()

	node := createNode(
		t,
		log,
		hub,
		collector,
		rootSnapshot,
	)

	<-node.Ready()
	<-node.Done()

	/*
				create a node
				send a transaction
				  verify it is included in a proposal on epoch1 cluster chain
			      by polling cluster state
				build epoch
		  		complete epoch (start epoch 2)
				  send a epoch 1 transaction
				    verify it is included in a proposal on epoch1 cluster chain
				  send a epoch 2 transaction
				    verify it is included in a proposal on epoch2 cluster chain

	*/
}
