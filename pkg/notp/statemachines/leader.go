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
	"fmt"

	notpsmpackets "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// NewLeaderStateMachine creates and configures a new leader state machine for the given operation.
func NewLeaderStateMachine(smType StateMachineType, hostHandler HostHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
	if smType == "" {
		smType = DefaultStateMachineType
	}
	stateMachine, err := NewStateMachine(smType, LeaderAdvertiseState, hostHandler, transportLayer)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create leader state machine: %w", err)
	}
	return stateMachine, nil
}

// LeaderAdvertiseState handles the advertisement phase in the protocol.
func LeaderAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	switch runtime.GetStateMachineType() {
	case PushStateMachineType:
		return false, nil, fmt.Errorf("notp: not implemented operation type: %s", runtime.GetStateMachineType())
	case PullStateMachineType:
		return false, notifyCurrentState, nil
	}
	return false, nil, fmt.Errorf("notp: unknown operation type: %s", runtime.GetStateMachineType())
}

// notifyCurrentState state to notify the current state.
func notifyCurrentState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.NotifyCurrentState)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondCurrentState, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle submit negotiation request packet: %w", err)
	}
	return false, FinalState, nil
}
