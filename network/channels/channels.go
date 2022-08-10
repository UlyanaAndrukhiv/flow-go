// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package channels

import (
	"fmt"
	"strings"

	"github.com/onflow/flow-go/model/flow"
)

// init is called first time this package is imported.
// It creates and initializes channelRoleMap and clusterChannelPrefixRoleMap.
func init() {
	initializeChannelRoleMap()
}

// channelRoleMap keeps a map between channels and the list of flow roles involved in them.
var channelRoleMap map[Channel]flow.RoleList

// clusterChannelPrefixRoleMap keeps a map between cluster channel prefixes and the list of flow roles involved in them.
var clusterChannelPrefixRoleMap map[string]flow.RoleList

// RolesByChannel returns list of flow roles involved in the channel.
// If the given channel is a public channel, the returned list will
// contain all roles.
func RolesByChannel(channel Channel) (flow.RoleList, bool) {
	if IsClusterChannel(channel) {
		return ClusterChannelRoles(channel), true
	}
	if PublicChannels().Contains(channel) {
		return flow.Roles(), true
	}
	roles, ok := channelRoleMap[channel]
	return roles, ok
}

// ChannelExists returns true if the channel exists.
func ChannelExists(channel Channel) bool {
	if _, ok := RolesByChannel(channel); ok {
		return true
	}

	return false
}

// ChannelsByRole returns a list of all channels the role subscribes to (except cluster-based channels and public channels).
func ChannelsByRole(role flow.Role) ChannelList {
	channels := make(ChannelList, 0)
	for channel, roles := range channelRoleMap {
		if roles.Contains(role) {
			channels = append(channels, channel)
		}
	}

	return channels
}

// UniqueChannels returns list of non-cluster channels with a unique RoleList accompanied
// with the list of all cluster channels.
// e.g. if channel X and Y both are non-cluster channels and have role IDs [A,B,C] then only one of them will be in the returned list.
func UniqueChannels(channels ChannelList) ChannelList {
	// uniques keeps the set of unique channels based on their RoleList.
	uniques := make(ChannelList, 0)
	// added keeps track of channels added to uniques for deduplication.
	added := make(map[flow.Identifier]struct{})

	// a channel is added to uniques if it is either a
	// cluster channel, or no non-cluster channel with the same set of roles
	// has already been added to uniques.
	// We use identifier of RoleList to determine its uniqueness.
	for _, channel := range channels {
		// non-cluster channel deduplicated based identifier of role list
		if !IsClusterChannel(channel) {
			id := channelRoleMap[channel].ID()
			if _, ok := added[id]; ok {
				// a channel with same RoleList already added, hence skips
				continue
			}
			added[id] = struct{}{}
		}

		uniques = append(uniques, channel)
	}

	return uniques
}

// Channels returns all channels that nodes of any role have subscribed to (except cluster-based channels).
func Channels() ChannelList {
	channels := make(ChannelList, 0)
	for channel := range channelRoleMap {
		channels = append(channels, channel)
	}
	channels = append(channels, PublicChannels()...)

	return channels
}

// PublicChannels returns all channels that are used on the public network.
func PublicChannels() ChannelList {
	return ChannelList{
		PublicSyncCommittee,
		PublicReceiveBlocks,
	}
}

// channels
const (

	// Channels used for testing
	TestNetworkChannel = Channel("test-network")
	TestMetricsChannel = Channel("test-metrics")

	// Channels for consensus protocols
	ConsensusCommittee     = Channel("consensus-committee")
	ConsensusClusterPrefix = "consensus-cluster" // dynamic channel, use ConsensusCluster function

	// Channels for protocols actively synchronizing state across nodes
	SyncCommittee     = Channel("sync-committee")
	SyncClusterPrefix = "sync-cluster" // dynamic channel, use SyncCluster function
	SyncExecution     = Channel("sync-execution")

	// Channels for dkg communication
	DKGCommittee = "dkg-committee"

	// Channels for actively pushing entities to subscribers
	PushTransactions = Channel("push-transactions")
	PushGuarantees   = Channel("push-guarantees")
	PushBlocks       = Channel("push-blocks")
	PushReceipts     = Channel("push-receipts")
	PushApprovals    = Channel("push-approvals")

	// Channels for actively requesting missing entities
	RequestCollections       = Channel("request-collections")
	RequestChunks            = Channel("request-chunks")
	RequestReceiptsByBlockID = Channel("request-receipts-by-block-id")
	RequestApprovalsByChunk  = Channel("request-approvals-by-chunk")

	// Channel aliases to make the code more readable / more robust to errors
	ReceiveTransactions = PushTransactions
	ReceiveGuarantees   = PushGuarantees
	ReceiveBlocks       = PushBlocks
	ReceiveReceipts     = PushReceipts
	ReceiveApprovals    = PushApprovals

	ProvideCollections       = RequestCollections
	ProvideChunks            = RequestChunks
	ProvideReceiptsByBlockID = RequestReceiptsByBlockID
	ProvideApprovalsByChunk  = RequestApprovalsByChunk

	// Public network channels
	PublicPushBlocks    = Channel("public-push-blocks")
	PublicReceiveBlocks = PublicPushBlocks
	PublicSyncCommittee = Channel("public-sync-committee")

	// Execution data service
	ExecutionDataService = Channel("execution-data-service")
)

// initializeChannelRoleMap initializes an instance of channelRoleMap and populates it
// with the channels and their corresponding list of authorized roles.
// Note: Please update this map, if a new channel is defined or a the roles subscribing to a channel have changed
func initializeChannelRoleMap() {
	channelRoleMap = make(map[Channel]flow.RoleList)

	// Channels for test
	channelRoleMap[TestNetworkChannel] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus, flow.RoleExecution,
		flow.RoleVerification, flow.RoleAccess}
	channelRoleMap[TestMetricsChannel] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus, flow.RoleExecution,
		flow.RoleVerification, flow.RoleAccess}

	// Channels for consensus protocols
	channelRoleMap[ConsensusCommittee] = flow.RoleList{flow.RoleConsensus}

	// Channels for protocols actively synchronizing state across nodes
	channelRoleMap[SyncCommittee] = flow.Roles()
	channelRoleMap[SyncExecution] = flow.RoleList{flow.RoleExecution}

	// Channels for DKG communication
	channelRoleMap[DKGCommittee] = flow.RoleList{flow.RoleConsensus}

	// Channels for actively pushing entities to subscribers
	channelRoleMap[PushTransactions] = flow.RoleList{flow.RoleCollection}
	channelRoleMap[PushGuarantees] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus}
	channelRoleMap[PushBlocks] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus, flow.RoleExecution,
		flow.RoleVerification, flow.RoleAccess}
	channelRoleMap[PushReceipts] = flow.RoleList{flow.RoleConsensus, flow.RoleExecution, flow.RoleVerification,
		flow.RoleAccess}
	channelRoleMap[PushApprovals] = flow.RoleList{flow.RoleConsensus, flow.RoleVerification}

	// Channels for actively requesting missing entities
	channelRoleMap[RequestCollections] = flow.RoleList{flow.RoleCollection, flow.RoleExecution, flow.RoleAccess}
	channelRoleMap[RequestChunks] = flow.RoleList{flow.RoleExecution, flow.RoleVerification}
	channelRoleMap[RequestReceiptsByBlockID] = flow.RoleList{flow.RoleConsensus, flow.RoleExecution}
	channelRoleMap[RequestApprovalsByChunk] = flow.RoleList{flow.RoleConsensus, flow.RoleVerification}

	// Channel aliases to make the code more readable / more robust to errors
	channelRoleMap[ReceiveGuarantees] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus}
	channelRoleMap[ReceiveBlocks] = flow.RoleList{flow.RoleCollection, flow.RoleConsensus, flow.RoleExecution,
		flow.RoleVerification, flow.RoleAccess}
	channelRoleMap[ReceiveReceipts] = flow.RoleList{flow.RoleConsensus, flow.RoleExecution, flow.RoleVerification,
		flow.RoleAccess}
	channelRoleMap[ReceiveApprovals] = flow.RoleList{flow.RoleConsensus, flow.RoleVerification}

	channelRoleMap[ProvideCollections] = flow.RoleList{flow.RoleCollection, flow.RoleExecution, flow.RoleAccess}
	channelRoleMap[ProvideChunks] = flow.RoleList{flow.RoleExecution, flow.RoleVerification}
	channelRoleMap[ProvideReceiptsByBlockID] = flow.RoleList{flow.RoleConsensus, flow.RoleExecution}
	channelRoleMap[ProvideApprovalsByChunk] = flow.RoleList{flow.RoleConsensus, flow.RoleVerification}

	clusterChannelPrefixRoleMap = make(map[string]flow.RoleList)

	clusterChannelPrefixRoleMap[SyncClusterPrefix] = flow.RoleList{flow.RoleCollection}
	clusterChannelPrefixRoleMap[ConsensusClusterPrefix] = flow.RoleList{flow.RoleCollection}
}

// ClusterChannelRoles returns the list of roles that are involved in the given cluster-based channel.
func ClusterChannelRoles(clusterChannel Channel) flow.RoleList {
	if prefix, ok := ClusterChannelPrefix(clusterChannel); ok {
		return clusterChannelPrefixRoleMap[prefix]
	}

	return flow.RoleList{}
}

// ClusterChannelPrefix returns the cluster channel prefix and true if clusterChannel exists inclusterChannelPrefixRoleMap
func ClusterChannelPrefix(clusterChannel Channel) (string, bool) {
	for prefix := range clusterChannelPrefixRoleMap {
		if strings.HasPrefix(clusterChannel.String(), prefix) {
			return prefix, true
		}
	}

	return "", false
}

// IsClusterChannel returns true if channel is cluster-based.
// Currently, only collection nodes are involved in a cluster-based channels.
func IsClusterChannel(channel Channel) bool {
	_, ok := ClusterChannelPrefix(channel)
	return ok
}

// TopicFromChannel returns the unique LibP2P topic form the channel.
// The channel is made up of name string suffixed with root block id.
// The root block id is used to prevent cross talks between nodes on different sporks.
func TopicFromChannel(channel Channel, rootBlockID flow.Identifier) Topic {
	// skip root block suffix, if this is a cluster specific channel. A cluster specific channel is inherently
	// unique for each epoch
	if IsClusterChannel(channel) {
		return Topic(channel)
	}
	return Topic(fmt.Sprintf("%s/%s", string(channel), rootBlockID.String()))
}

func ChannelFromTopic(topic Topic) (Channel, bool) {
	if IsClusterChannel(Channel(topic)) {
		return Channel(topic), true
	}

	if index := strings.LastIndex(topic.String(), "/"); index != -1 {
		return Channel(topic[:index]), true
	}

	return "", false
}

// ConsensusCluster returns a dynamic cluster consensus channel based on
// the chain ID of the cluster in question.
func ConsensusCluster(clusterID flow.ChainID) Channel {
	return Channel(fmt.Sprintf("%s-%s", ConsensusClusterPrefix, clusterID))
}

// SyncCluster returns a dynamic cluster sync channel based on the chain
// ID of the cluster in question.
func SyncCluster(clusterID flow.ChainID) Channel {
	return Channel(fmt.Sprintf("%s-%s", SyncClusterPrefix, clusterID))
}