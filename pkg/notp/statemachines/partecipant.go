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

// FlowType represents the type of operation that the NOTP protocol is performing.
type FlowType string

const (
	// PushFlowType represents the push state machine type.
	PushFlowType FlowType = "push"
	// PullFlowType represents the pull state machine type.
	PullFlowType FlowType = "pull"
	// DefaultFlowType represents the default operation type.
	DefaultFlowType FlowType = PushFlowType
)

// notifyProtocol state to notify the protocol.
func startFlow(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.StartFlowMessage, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle start flow packet: %w", err)
	}
	_, _, _, err = receiveAndHandleStatePacket(runtime, notpsmpackets.ActionResponseMessage)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle action response packet: %w", err)
	}
	return false, subscriberNegotiationState, nil
}

// processStartFlow state to process the start flow.
func processStartFlow(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.StartFlowMessage)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle start flow packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.ActionResponseMessage, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle action response packet: %w", err)
	}
	return false, FinalState, nil
}

// notifyProtocol state to notify the protocol.
func notifyProtocol(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RequestCurrentObjectsStateMessage, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle request current state packet: %w", err)
	}
	_, _, _, err = receiveAndHandleStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	return false, subscriberNegotiationState, nil
}

// requestObjectsState state to request the current state.
func requestObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RequestCurrentObjectsStateMessage, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle request current state packet: %w", err)
	}
	_, _, _, err = receiveAndHandleStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	return false, subscriberNegotiationState, nil
}

// processRequestObjectsState state to process the request for the current state.
func processRequestObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, _, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.RequestCurrentObjectsStateMessage)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to receive and handle request current state packet: %w", err)
	}
	err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage, packetables)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	return false, FinalState, nil
}

// notifyObjectsState state to send the current state notification.
func notifyObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.NotifyCurrentObjectStatesMessage, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle notify current state packet: %w", err)
	}
	return false, publisherNegotiationState, nil
}

// processNotifyObjectsState state to process the current state notification.
func processNotifyObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type")
}

// submitNegotiationResponse state to submit negotiation response.
func subscriberNegotiationState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.NegotiationRequestMessage, nil)
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle submit negotiation request packet: %w", err)
	}
	return false, FinalState, nil
}

// submitNegotiationResponse state to submit negotiation response.
func publisherNegotiationState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type")
}

// publisherDataStreamState state to send data stream.
func publisherDataStreamState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type")
}

// subscriberDataStreamState state to receive data stream.
func subscriberDataStreamState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, nil, fmt.Errorf("notp: not implemented operation type")
}
