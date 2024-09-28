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
	"errors"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notpsmpackets "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// HandlerContext holds the context of the handler.
type HandlerContext struct {
	flow FlowType
}

// GetFlowType returns the flow type of the handler context.
func (h *HandlerContext) GetFlowType() FlowType {
	return h.flow
}

// PacketCreatorFunc is a function that creates a packet.
type PacketCreatorFunc func(*notpsmpackets.StatePacket) notppackets.Packetable

// HostHandler defines a function type for handling packet.
type HostHandler func(*HandlerContext, *notpsmpackets.StatePacket, []notppackets.Packetable) (bool, []notppackets.Packetable, error)

// StateTransitionFunc defines a function responsible for transitioning to the next state in the state machine.
type StateTransitionFunc func(runtime *StateMachineRuntimeContext) (isFinal bool, nextState StateTransitionFunc, err error)

// InitialState defines the initial state of the state machine.
func InitialState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, runtime.initialState, nil
}

// FinalState defines the final state of the state machine.
func FinalState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return true, nil, nil
}

// StateMachineRuntimeContext holds the runtime context of the state machine.
type StateMachineRuntimeContext struct {
	flow		   FlowType
	transportLayer *notptransport.TransportLayer
	initialState   StateTransitionFunc
	hostHandler    HostHandler
}

// WithFlow returns the state machine runtime context with the flow type.
func (t *StateMachineRuntimeContext) WithFlow(flowType FlowType) *StateMachineRuntimeContext {
	return &StateMachineRuntimeContext{
		flow:           flowType,
		transportLayer: t.transportLayer,
		initialState:   t.initialState,
		hostHandler:    t.hostHandler,
	}
}

// GetFlowType returns the flow type of the state machine.
func (t *StateMachineRuntimeContext) GetFlowType() FlowType {
	return t.flow
}

// Send sends a packet through the transport layer.
func (t *StateMachineRuntimeContext) Send(packetable notppackets.Packetable) error {
	return t.SendStream([]notppackets.Packetable{packetable})
}

// SendStream sends a packets through the transport layer.
func (t *StateMachineRuntimeContext) SendStream(packetables []notppackets.Packetable) error {
	return t.transportLayer.TransmitPacket(packetables)
}

// Receive retrieves a packet from the transport layer.
func (t *StateMachineRuntimeContext) Receive() (notppackets.Packetable, error) {
	packets, err := t.ReceiveStream()
	if err != nil {
		return nil, err
	}
	if len(packets) == 0 {
		return nil, errors.New("notp: received a nil packet")
	} else if len(packets) > 1 {
		return nil, errors.New("notp: received more than one packet")
	}
	return packets[0], nil
}

// ReceiveStream retrieves packets from the transport layer.
func (t *StateMachineRuntimeContext) ReceiveStream() ([]notppackets.Packetable, error) {
	return t.transportLayer.ReceivePacket()
}

// Handle handles the packet for the state machine.
func (t *StateMachineRuntimeContext) Handle(handlerCtx *HandlerContext, statePacket *notpsmpackets.StatePacket) (bool, []notppackets.Packetable, error) {
	return t.HandleStream(handlerCtx, statePacket, nil)
}

// HandleStream handles a packet stream for the state machine.
func (t *StateMachineRuntimeContext) HandleStream(handlerCtx *HandlerContext, statePacket *notpsmpackets.StatePacket, packetables []notppackets.Packetable) (bool, []notppackets.Packetable, error) {
	if packetables == nil {
		packetables = []notppackets.Packetable{}
	}
	return t.hostHandler(handlerCtx, statePacket, packetables)
}

// StateMachine orchestrates the execution of state transitions.
type StateMachine struct {
	runtime *StateMachineRuntimeContext
}

// Run starts and runs the state machine through its states until termination.
func (m *StateMachine) Run() error {
	state := m.runtime.initialState
	for state != nil {
		isFinal, nextState, err := state(m.runtime)
		if err != nil {
			return err
		}
		if isFinal {
			break
		}
		state = nextState
	}
	return nil
}

// NewStateMachine creates and initializes a new state machine with the given initial state and transport layer.
func NewStateMachine(initialState StateTransitionFunc, hostHandler HostHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
	if initialState == nil {
		return nil, errors.New("notp: initial state cannot be nil")
	}
	if hostHandler == nil {
		return nil, errors.New("notp: decision handler cannot be nil")
	}
	if transportLayer == nil {
		return nil, errors.New("notp: transport layer cannot be nil")
	}
	return &StateMachine{
		runtime: &StateMachineRuntimeContext{
			transportLayer: transportLayer,
			initialState:   initialState,
			hostHandler:    hostHandler,
		},
	}, nil
}
