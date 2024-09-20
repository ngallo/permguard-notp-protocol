// Copyright 2024 Nitro Agility S.r.l.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package statemachines

import (
	"testing"

	"github.com/stretchr/testify/assert"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// stateMachinesInfo represents the state machines and their respective packet logs.
type stateMachinesInfo struct {
	follower        *StateMachine
	followerSent    []notppackets.Packet
	followerReceived []notppackets.Packet

	leader          *StateMachine
	leaderSent      []notppackets.Packet
	leaderReceived  []notppackets.Packet
}

// buildCommitStateMachines initializes and returns the follower and leader state machines.
func buildCommitStateMachines(assert *assert.Assertions, operationType OperationType, followerHandler DecisionHandler, leaderHandler DecisionHandler) *stateMachinesInfo {
	sMInfo := &stateMachinesInfo{
		followerSent:     []notppackets.Packet{},
		followerReceived: []notppackets.Packet{},
		leaderSent:       []notppackets.Packet{},
		leaderReceived:   []notppackets.Packet{},
	}

	onFollowerSent := func(packet *notppackets.Packet) {
		sMInfo.followerSent = append(sMInfo.followerSent, *packet)
	}
	onFollowerReceived := func(packet *notppackets.Packet) {
		sMInfo.followerReceived = append(sMInfo.followerReceived, *packet)
	}

	onLeaderSent := func(packet *notppackets.Packet) {
		sMInfo.leaderSent = append(sMInfo.leaderSent, *packet)
	}
	onLeaderReceived := func(packet *notppackets.Packet) {
		sMInfo.leaderReceived = append(sMInfo.leaderReceived, *packet)
	}

	followerStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the follower transport stream")
	leaderStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the leader transport stream")

	followerPacketLogger, err := notptransport.NewPacketInspector(onFollowerSent, onFollowerReceived)
	assert.Nil(err, "Failed to initialize the follower packet logger")
	followerTransport, err := notptransport.NewTransportLayer(leaderStream.TransmitPacket, followerStream.ReceivePacket, followerPacketLogger)
	assert.Nil(err, "Failed to initialize the follower transport layer")

	leaderPacketLogger, err := notptransport.NewPacketInspector(onLeaderSent, onLeaderReceived)
	assert.Nil(err, "Failed to initialize the leader packet logger")
	leaderTransport, err := notptransport.NewTransportLayer(followerStream.TransmitPacket, leaderStream.ReceivePacket, leaderPacketLogger)
	assert.Nil(err, "Failed to initialize the leader transport layer")

	followerSMachine, err := NewFollowerStateMachine(operationType, followerHandler, followerTransport)
	assert.Nil(err, "Failed to initialize the follower state machine")
	sMInfo.follower = followerSMachine

	leaderSMachine, err := NewLeaderStateMachine(operationType, leaderHandler, leaderTransport)
	assert.Nil(err, "Failed to initialize the leader state machine")
	sMInfo.leader = leaderSMachine

	return sMInfo
}

// TestPullProtocolExecution verifies the state machine execution for both follower and leader in the context of a pull operation.
func TestPullProtocolExecution(t *testing.T) {
	assert := assert.New(t)

	followerHandler := func(packet *notppackets.Packet) (*notppackets.Packet, error) {
		return packet, nil
	}
	leaderHandler := func(packet *notppackets.Packet) (*notppackets.Packet, error) {
		return packet, nil
	}
	sMInfo := buildCommitStateMachines(assert, PullOperation, followerHandler, leaderHandler)

	err := sMInfo.follower.Run()
	assert.Nil(err, "Failed to run the follower state machine")

	err = sMInfo.leader.Run()
	assert.Nil(err, "Failed to run the leader state machine")
}
