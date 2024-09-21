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
	"fmt"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
	notpsmpackets "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines/packets"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// HandlerContext holds the context of the handler.
type HandlerContext struct {
	stateMachineType StateMachineType
	isLeader		 bool
}

// GetStateMachineType returns the state machine type of the handler.
func (h *HandlerContext) GetStateMachineType() StateMachineType {
	return h.stateMachineType
}

// IsLeader returns whether the handler is a leader.
func (h *HandlerContext) IsLeader() bool {
	return h.isLeader
}

// createBasePacketWithContext creates and initializes a new base packet with the given handler context.
func createBasePacketWithContext(smType StateMachineType, isLeader bool, opCode, algoCode uint16) (*notpsmpackets.BasePacket, *HandlerContext, error) {
	handlerCtx := &HandlerContext{
		stateMachineType: smType,
		isLeader: isLeader,
	}
	packet := &notpsmpackets.BasePacket{
		OperationCode: opCode,
		AlgorithmCode: algoCode,
		ErrorCode: 0,
	}
	return packet, handlerCtx, nil
}

// PacketCreatorFunc is a function that creates a packet.
type PacketCreatorFunc func(*notpsmpackets.BasePacket) notppackets.Packetable

// CrateAndHandlePacket creates and handles a packet.
func CrateAndHandlePacket(runtime *StateMachineRuntimeContext, smType StateMachineType, isLeader bool, opCode, algoCode uint16, packetCreator PacketCreatorFunc) (notppackets.Packetable, []notppackets.Packetable, error) {
	packet, handlerCtx, err := createBasePacketWithContext(smType, isLeader, opCode, algoCode)
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to create base packet with context: %w", err)
	}
	createdPacket := packetCreator(packet)
	packetables, err := runtime.Handle(handlerCtx, createdPacket)
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to handle created packet: %w", err)
	}
	return createdPacket, packetables, nil
}

// ReceiveAndHandleHeadStream receives and handles a head stream.
func ReceiveAndHandleHeadStream(runtime *StateMachineRuntimeContext, smType StateMachineType, isLeader bool, target notppackets.Packetable) ([]notppackets.Packetable, error) {
	handlerCtx := &HandlerContext{
		stateMachineType: PullStateMachineType,
		isLeader:         true,
	}
	packetsStream, err := runtime.ReceiveStream()
	if err != nil {
		return nil, fmt.Errorf("notp: failed to receive packets: %w", err)
	}
	err = notppackets.ConvertPacketable(packetsStream[0], target)
	if err != nil {
		return nil, fmt.Errorf("notp: failed to convert packetable: %w", err)
	}
	packetables, err := runtime.HandleStream(handlerCtx, packetsStream[1:])
	if err != nil {
		return nil, fmt.Errorf("notp: failed to handle packet stream: %w", err)
	}
	return packetables, nil
}

// HostHandler defines a function type for handling packet.
type HostHandler func(*HandlerContext, []notppackets.Packetable) ([]notppackets.Packetable, error)

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
	operation      StateMachineType
	transportLayer *notptransport.TransportLayer
	initialState   StateTransitionFunc
	hostHandler    HostHandler
}

// GetOperation returns the operation type of the state machine.
func (t *StateMachineRuntimeContext) GetOperation() StateMachineType {
	return t.operation
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
func (t *StateMachineRuntimeContext) Handle(handlerCtx *HandlerContext, packetable notppackets.Packetable) ([]notppackets.Packetable, error) {
	return t.HandleStream(handlerCtx, []notppackets.Packetable{packetable})
}

// HandleStream handles a packet stream for the state machine.
func (t *StateMachineRuntimeContext) HandleStream(handlerCtx *HandlerContext, packetables []notppackets.Packetable) ([]notppackets.Packetable, error) {
	return t.hostHandler(handlerCtx, packetables)
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
func NewStateMachine(operation StateMachineType, initialState StateTransitionFunc, hostHandler HostHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
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
			operation:      operation,
			transportLayer: transportLayer,
			initialState:   initialState,
			hostHandler:    hostHandler,
		},
	}, nil
}
