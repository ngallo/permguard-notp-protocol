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

// TestClientServerStateMachineExecution verifies the state machine execution for both client and server.
func TestClientServerStateMachineExecution(t *testing.T) {
	assert := assert.New(t)

	// Initialize in-memory streams for client and server communication
	clientStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the client transport stream")
	serverStream, err := notptransport.NewInMemoryStream()
	assert.Nil(err, "Failed to initialize the server transport stream")

	// Create transport layers for both client and server
	clientTransport, err := notptransport.NewTransportLayer(serverStream.TransmitPacket, clientStream.ReceivePacket, nil)
	assert.Nil(err, "Failed to initialize the client transport layer")
	serverTransport, err := notptransport.NewTransportLayer(clientStream.TransmitPacket, serverStream.ReceivePacket, nil)
	assert.Nil(err, "Failed to initialize the server transport layer")

	decisionHandler := func(packet *notppackets.Packet) (*notppackets.Packet, error) {
		return packet, nil
	}

	// Initialize and run client state machine
	clientSMachine, err := NewStateMachine(FinalState, decisionHandler, clientTransport)
	assert.Nil(err, "Failed to initialize the client state machine")
	err = clientSMachine.Run()
	assert.Nil(err, "Failed to run the client state machine")

	// Initialize and run server state machine
	serverSMachine, err := NewStateMachine(FinalState, decisionHandler, serverTransport)
	assert.Nil(err, "Failed to initialize the server state machine")
	err = serverSMachine.Run()
	assert.Nil(err, "Failed to run the server state machine")
}
