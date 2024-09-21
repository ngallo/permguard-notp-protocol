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

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notpsmpackets "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// NewLeaderStateMachine creates and configures a new leader state machine for the given operation.
func NewLeaderStateMachine(operation StateMachineType, hostHandler HostHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
	if operation == "" {
		operation = DefaultOperation
	}
	stateMachine, err := NewStateMachine(operation, LeaderAdvertiseState, hostHandler, transportLayer)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create leader state machine: %w", err)
	}
	return stateMachine, nil
}


// leaderPullAdvertiseState handles the pull advertisement phase in the protocol.
func leaderPullAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	// Receive the advertisement packet and its stream
	var advPacket notpsmpackets.AdvertisementPacket
	_, err := ReceiveHeadStream(runtime, &advPacket)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to convert packetable: %w", err)
	}

	// Send the advertisement packets stream and transition to the next state.
	return false, LeaderNegotiateState, nil
}

// leaderPushAdvertiseState handles the push advertisement phase in the protocol.
func leaderPushAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, LeaderNegotiateState, nil
}

// LeaderAdvertiseState handles the advertisement phase in the protocol.
func LeaderAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	if runtime.GetOperation() == PullStateMachineType {
		return leaderPullAdvertiseState(runtime)
	} else if runtime.GetOperation() == PushStateMachineType {
		return leaderPushAdvertiseState(runtime)
	}
	return false, nil, fmt.Errorf("notp: invalid operation type: %s", runtime.GetOperation())
}

// LeaderNegotiateState manages the negotiation phase in the protocol.
func LeaderNegotiateState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, LeaderExchangeState, nil
}

// LeaderExchangeState governs the exchange phase in the protocol.
func LeaderExchangeState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, FinalState, nil
}
