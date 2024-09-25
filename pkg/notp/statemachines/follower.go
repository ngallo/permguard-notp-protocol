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

// NewFollowerStateMachine creates and configures a new follower state machine for the given operation.
func NewFollowerStateMachine(operation StateMachineType, hostHandler HostHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
	if operation == "" {
		operation = DefaultOperation
	}
	stateMachine, err := NewStateMachine(operation, FollowerAdvertiseState, hostHandler, transportLayer)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create follower state machine: %w", err)
	}
	return stateMachine, nil
}

// FollowerAdvertiseState handles the advertisement phase in the protocol.
func FollowerAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	switch runtime.GetOperation() {
	case PushStateMachineType:
		return false, requestCurrentState, nil
	case PullStateMachineType:
		return false, requestCurrentState, nil
	}
	return false, nil, fmt.Errorf("notp: unknown operation type: %s", runtime.GetOperation())
}

// requestCurrentState state to request the current state.
func requestCurrentState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, runtime.GetOperation(), false, notpsmpackets.NotifyCurrentState, notpsmpackets.AlgoNone, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle request current state packet: %w", err)
	}
	return false, requestCurrentState, nil
}

// submitNegotiationRequest state to submit negotiation request.
func submitNegotiationRequest(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, runtime.GetOperation(), false, notpsmpackets.RespondCurrentState)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, runtime.GetOperation(), false, notpsmpackets.SubmitNegotiationRequest, notpsmpackets.AlgoNone, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle submit negotiation request packet: %w", err)
	}
	return false, FinalState, nil
}
