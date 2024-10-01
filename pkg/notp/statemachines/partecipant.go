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
)

// FlowType represents the type of operation that the NOTP protocol is performing.
type FlowType uint64

const (
	// UnknownFlowType represents an unknown state machine type.
	UnknownFlowType FlowType = 0
	// PushFlowType represents the push state machine type.
	PushFlowType FlowType = 1
	// PullFlowType represents the pull state machine type.
	PullFlowType FlowType = 2
	// DefaultFlowType represents the default operation type.
	DefaultFlowType FlowType = PushFlowType

	// StartFlowStateID represents the state ID for the start flow state.
	StartFlowStateID = uint16(10)
	// ProcessStartFlowStateID represents the state ID for the process start flow state.
	ProcessStartFlowStateID = uint16(11)
	// RequestObjectsStateID represents the state ID for the request objects state.
	RequestObjectsStateID = uint16(13)
	// ProcessRequestObjectsStateID represents the state ID for the process request objects state.
	ProcessRequestObjectsStateID = uint16(14)
	// NotifyObjectsStateID represents the state ID for the notify objects state.
	NotifyObjectsStateID = uint16(15)
	// ProcessNotifyObjectsStateID represents the state ID for the process notify objects state.
	ProcessNotifyObjectsStateID = uint16(16)
	// SubscriberNegotiationStateID represents the state ID for the subscriber negotiation state.
	SubscriberNegotiationStateID = uint16(17)
	// SubscriberDataStreamStateID represents the state ID for the subscriber data stream state.
	SubscriberDataStreamStateID = uint16(18)
	// PublisherNegotiationStateID represents the state ID for the publisher negotiation state.
	PublisherNegotiationStateID = uint16(19)
	// PublisherDataStreamStateID represents the state ID for the publisher data stream state.
	PublisherDataStreamStateID = uint16(20)
)

// defaultStateMap represents the default state map for the state machine.
var defaultStateMap = map[uint16]StateTransitionFunc{
	InitialStateID:               InitialState,
	FinalStateID:                 FinalState,
	StartFlowStateID:             startFlowState,
	ProcessStartFlowStateID:      processStartFlowState,
	RequestObjectsStateID:        requestObjectsState,
	ProcessRequestObjectsStateID: processRequestObjectsState,
	NotifyObjectsStateID:         notifyObjectsState,
	ProcessNotifyObjectsStateID:  processNotifyObjectsState,
	PublisherNegotiationStateID:  publisherNegotiationState,
	PublisherDataStreamStateID:   publisherDataStreamState,
	SubscriberNegotiationStateID: subscriberNegotiationState,
	SubscriberDataStreamStateID:  subscriberDataStreamState,
}

// startFlowState state to start the flow.
func startFlowState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, err := createAndHandleAndStreamStatePacketWithValue(runtime, notpsmpackets.StartFlowMessage, uint64(runtime.flowType), nil)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle start flow packet: %w", err)
	}
	statePacket, _, err := receiveAndHandleStatePacket(runtime, notpsmpackets.ActionResponseMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle action response packet: %w", err)
	}
	if !statePacket.HasAck() {
		return nil, fmt.Errorf("notp: failed to receive ack in action response packet")
	}
	var stateID uint16
	switch runtime.GetFlowType() {
	case PushFlowType:
		stateID = NotifyObjectsStateID
	case PullFlowType:
		stateID = RequestObjectsStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// processStartFlowState state to process the start flow.
func processStartFlowState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	statePacket, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.StartFlowMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle start flow packet: %w", err)
	}
	messageValue := notppackets.CombineUint32toUint64(notpsmpackets.AcknowledgedValue, notpsmpackets.UnknownValue)
	_, err = createAndHandleAndStreamStatePacketWithValue(runtime, notpsmpackets.ActionResponseMessage, messageValue, packetables)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle action response packet: %w", err)
	}
	flowtype := FlowType(statePacket.MessageValue)
	runtime = runtime.WithFlow(flowtype)
	var stateID uint16
	switch runtime.GetFlowType() {
	case PushFlowType:
		stateID = ProcessNotifyObjectsStateID
	case PullFlowType:
		stateID = ProcessRequestObjectsStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// requestObjectsState state to request the current state.
func requestObjectsState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RequestCurrentObjectsStateMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle request current state packet: %w", err)
	}
	_, _, err = receiveAndHandleStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	var stateID uint16
	switch runtime.GetFlowType() {
	case PullFlowType:
		stateID = SubscriberNegotiationStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// processRequestObjectsState state to process the request for the current state.
func processRequestObjectsState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.RequestCurrentObjectsStateMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle request current state packet: %w", err)
	}
	_, err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage, packetables)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	var stateID uint16
	switch runtime.GetFlowType() {
	case PullFlowType:
		stateID = PublisherNegotiationStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// notifyObjectsState state to send the current state notification.
func notifyObjectsState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.NotifyCurrentObjectStatesMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle notify current state packet: %w", err)
	}
	_, _, err = receiveAndHandleStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	var stateID uint16
	switch runtime.GetFlowType() {
	case PushFlowType:
		stateID = PublisherNegotiationStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// processNotifyObjectsState state to process the current state notification.
func processNotifyObjectsState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.NotifyCurrentObjectStatesMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle notify current state packet: %w", err)
	}
	_, err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondCurrentStateMessage, packetables)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	var stateID uint16
	switch runtime.GetFlowType() {
	case PushFlowType:
		stateID = SubscriberNegotiationStateID
	default:
		return nil, fmt.Errorf("notp: unknown flow type")
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// submitNegotiationResponse state to submit negotiation response.
func subscriberNegotiationState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.NegotiationRequestMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle notify current state packet: %w", err)
	}
	statePacket, _, err := receiveAndHandleStatePacket(runtime, notpsmpackets.RespondNegotiationRequestMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle respond current state packet: %w", err)
	}
	if !statePacket.HasAck() {
		return nil, fmt.Errorf("notp: failed to receive ack in respond negotiation request packet")
	}
	stateID := SubscriberDataStreamStateID
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// submitNegotiationResponse state to submit negotiation response.
func publisherNegotiationState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, packetables, err := receiveAndHandleStatePacket(runtime, notpsmpackets.NegotiationRequestMessage)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive and handle notify current state packet: %w", err)
	}
	_, err = createAndHandleAndStreamStatePacket(runtime, notpsmpackets.RespondNegotiationRequestMessage, packetables)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	stateID := PublisherDataStreamStateID
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: stateID,
	}, nil
}

// subscriberDataStreamState state to receive data stream.
func subscriberDataStreamState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	hasStream := true
	for hasStream {
		statePacket, _, err := receiveAndHandleStatePacket(runtime, notpsmpackets.ExchangeDataStreamMessage)
		if err != nil {
			return nil, fmt.Errorf("notp: failed to receive and handle exchange data stream packet: %w", err)
		}
		hasStream = statePacket.HasActiveDataStream()
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: FinalStateID,
	}, nil
}

// publisherDataStreamState state to send data stream.
func publisherDataStreamState(runtime *StateMachineRuntimeContext) (*StateTransitionInfo, error) {
	_, err := createAndHandleAndStreamStatePacket(runtime, notpsmpackets.ExchangeDataStreamMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to create and handle respond current state packet: %w", err)
	}
	return &StateTransitionInfo{
		Runtime: runtime,
		StateID: FinalStateID,
	}, nil
}
