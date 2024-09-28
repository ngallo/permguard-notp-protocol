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
)

// StateMachineType represents the type of operation that the NOTP protocol is performing.
type StateMachineType string

const (
	// PushStateMachineType represents the push state machine type.
	PushStateMachineType StateMachineType = "push"
	// PullStateMachineType represents the pull state machine type.
	PullStateMachineType StateMachineType = "pull"
	// DefaultStateMachineType represents the default operation type.
	DefaultStateMachineType StateMachineType = PushStateMachineType
)

// requestCurrentState state to request the current state.
func requestCurrentState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RequestCurrentState, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle request current state packet: %w", err)
	}
	return false, respondNegotiationRequest, nil
}

// notifyCurrentState state to notify the current state.
func notifyCurrentState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.NotifyCurrentState, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle notify current state packet: %w", err)
	}
	return false, handleNegotiationResponse, nil
}

// respondCurrentState state to respond to the current state.
func respondCurrentState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.RequestCurrentState)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle request current state packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondCurrentState, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	return false, FinalState, nil
}

// respondNegotiationRequest state to submit negotiation request.
func respondNegotiationRequest(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.RespondCurrentState)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.SubmitNegotiationRequest, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle submit negotiation request packet: %w", err)
	}
	return false, FinalState, nil
}

// submitNegotiationResponse state to submit negotiation response.
func handleNegotiationRequest(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type: %s", runtime.GetStateMachineType())
}

// submitNegotiationResponse state to submit negotiation response.
func handleNegotiationResponse(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type: %s", runtime.GetStateMachineType())
}

// sendDataStream state to send data stream.
func sendDataStream(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type: %s", runtime.GetStateMachineType())
}

// receiveDataStream state to receive data stream.
func receiveDataStream(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type: %s", runtime.GetStateMachineType())
}
