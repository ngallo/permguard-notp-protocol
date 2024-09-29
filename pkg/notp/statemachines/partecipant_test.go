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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notpsmpackets "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// stateMachinesInfo represents the state machines and their respective packet logs.
type stateMachinesInfo struct {
	follower         *StateMachine
	followerSent     []notppackets.Packet
	followerReceived []notppackets.Packet

	leader         *StateMachine
	leaderSent     []notppackets.Packet
	leaderReceived []notppackets.Packet
}

// buildCommitStateMachines initializes and returns the follower and leader state machines.
func buildCommitStateMachines(assert *assert.Assertions, followerHandler HostHandler, leaderHandler HostHandler) *stateMachinesInfo {
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

	followerStream, err := notptransport.NewInMemoryStream(5 * time.Second)
	assert.Nil(err, "Failed to initialize the follower transport stream")
	leaderStream, err := notptransport.NewInMemoryStream(5 * time.Second)
	assert.Nil(err, "Failed to initialize the leader transport stream")

	followerPacketLogger, err := notptransport.NewPacketInspector(onFollowerSent, onFollowerReceived)
	assert.Nil(err, "Failed to initialize the follower packet logger")
	followerTransport, err := notptransport.NewTransportLayer(leaderStream.TransmitPacket, followerStream.ReceivePacket, followerPacketLogger)
	assert.Nil(err, "Failed to initialize the follower transport layer")

	leaderPacketLogger, err := notptransport.NewPacketInspector(onLeaderSent, onLeaderReceived)
	assert.Nil(err, "Failed to initialize the leader packet logger")
	leaderTransport, err := notptransport.NewTransportLayer(followerStream.TransmitPacket, leaderStream.ReceivePacket, leaderPacketLogger)
	assert.Nil(err, "Failed to initialize the leader transport layer")

	followerSMachine, err := NewFollowerStateMachine(followerHandler, followerTransport)
	assert.Nil(err, "Failed to initialize the follower state machine")
	sMInfo.follower = followerSMachine

	leaderSMachine, err := NewLeaderStateMachine(leaderHandler, leaderTransport)
	assert.Nil(err, "Failed to initialize the leader state machine")
	sMInfo.leader = leaderSMachine

	return sMInfo
}

// TestPullProtocolExecution verifies the state machine execution for both follower and leader in the context of a pull operation.
func TestPullProtocolExecution(t *testing.T) {
	assert := assert.New(t)

	followerHandler := func(handlerCtx *HandlerContext, statePacket *notpsmpackets.StatePacket, packets []notppackets.Packetable) (bool, uint64, []notppackets.Packetable, uint16, error) {
		if !statePacket.HasAck() {
			return false, notpsmpackets.ActionRejected, packets, 0, nil
		}
		return false, notpsmpackets.ActionAcknowledged, packets, 0, nil
	}
	leaderHandler := func(handlerCtx *HandlerContext, statePacket *notpsmpackets.StatePacket, packets []notppackets.Packetable) (bool, uint64, []notppackets.Packetable, uint16, error) {
		return false, notpsmpackets.ActionAcknowledged, packets, 0, nil
	}
	sMInfo := buildCommitStateMachines(assert, followerHandler, leaderHandler)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := sMInfo.follower.Run(PullFlowType)
		assert.Nil(err, "Failed to run the follower state machine")
	}()

	go func() {
		defer wg.Done()
		err := sMInfo.leader.Run(UnknownFlowType)
		assert.Nil(err, "Failed to run the leader state machine")
	}()

	wg.Wait()

	assert.Len(sMInfo.followerSent, 2, "Follower sent packets")
	assert.Len(sMInfo.followerReceived, 1, "Follower received packets")
	assert.Len(sMInfo.leaderSent, 1, "Leader sent packets")
	assert.Len(sMInfo.leaderReceived, 1, "Leader received packets")
}
