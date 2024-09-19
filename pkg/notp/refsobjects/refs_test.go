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

package refsobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

type statesMachinesInfo struct {
	clientSent     []notppackets.Packet
	clientReceived []notppackets.Packet
	serverSent     []notppackets.Packet
	serverReceived []notppackets.Packet
	clientSMachine *ClientStateMachine
	serverSMachine *ServerStateMachine
}

// buildRefsObjectsStateMachines initializes the client and server state machines.
func buildRefsObjectsStateMachines(assert *assert.Assertions) *statesMachinesInfo {
	sMInfo := &statesMachinesInfo{
		clientSent:     []notppackets.Packet{},
		clientReceived: []notppackets.Packet{},
		serverSent:     []notppackets.Packet{},
		serverReceived: []notppackets.Packet{},
	}
	onClientSent := func(packet *notppackets.Packet) {
		sMInfo.clientSent = append(sMInfo.clientSent, *packet)
	}
	onClientReceived := func(packet *notppackets.Packet) {
		sMInfo.clientReceived = append(sMInfo.clientReceived, *packet)
	}
	clientStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the client transport stream")
	serverStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the server transport stream")

	clientPacketLogger, err := notptransport.NewPacketLogger(onClientSent, onClientReceived)
	assert.Nil(err, "Failed to initialize the packet logger")
	clientTransport, err := notptransport.NewTransportLayer(serverStream.TransmitPacket, clientStream.ReceivePacket, clientPacketLogger)
	assert.Nil(err, "Failed to initialize the client transport layer")
	onServerSent := func(packet *notppackets.Packet) {
		sMInfo.serverSent = append(sMInfo.serverSent, *packet)
	}
	onServerReceived := func(packet *notppackets.Packet) {
		sMInfo.serverReceived = append(sMInfo.serverReceived, *packet)
	}
	serverPacketLogger, err := notptransport.NewPacketLogger(onServerSent, onServerReceived)
	assert.Nil(err, "Failed to initialize the packet logger")
	serverTransport, err := notptransport.NewTransportLayer(clientStream.TransmitPacket, serverStream.ReceivePacket, serverPacketLogger)
	assert.Nil(err, "Failed to initialize the server transport layer")

	clientSMachine, err := NewClientStateMachine(clientTransport)
	assert.Nil(err, "Failed to initialize the client state machine")
	sMInfo.clientSMachine = clientSMachine
	serverSMachine, err := NewServerStateMachine(serverTransport)
	assert.Nil(err, "Failed to initialize the server state machine")
	sMInfo.serverSMachine = serverSMachine
	return sMInfo
}

// TestClientServerStateMachineExecution verifies the state machine execution for both client and server.
func TestClientServerStateMachineExecution(t *testing.T) {
	assert := assert.New(t)
	sMInfo := buildRefsObjectsStateMachines(assert)
	err := sMInfo.clientSMachine.Run()
	assert.Nil(err, "Failed to run the client state machine")
	sMInfo.serverSMachine.Run()
	assert.Nil(err, "Failed to run the server state machine")
}
