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
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// DecisionHandler defines a function type for handling packet dections.
type DecisionHandler func(*notppackets.Packet) (*notppackets.Packet, error)

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
	transportLayer *notptransport.TransportLayer
	initialState   StateTransitionFunc
	decisionHandler DecisionHandler
}

// TransmitPacket sends a packet through the transport layer.
func (t *StateMachineRuntimeContext) TransmitPacket(packet *notppackets.Packet) error {
	return t.transportLayer.TransmitPacket(packet)
}

// ReceivePacket retrieves a packet from the transport layer.
func (t *StateMachineRuntimeContext) ReceivePacket() (*notppackets.Packet, error) {
	return t.transportLayer.ReceivePacket()
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
func NewStateMachine(initialState StateTransitionFunc, decisionHandler DecisionHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
	if initialState == nil {
		return nil, errors.New("notp: initial state cannot be nil")
	}
	if decisionHandler == nil {
		return nil, errors.New("notp: decision handler cannot be nil")
	}
	if transportLayer == nil {
		return nil, errors.New("notp: transport layer cannot be nil")
	}
	return &StateMachine{
		runtime: &StateMachineRuntimeContext{
			transportLayer: transportLayer,
			initialState:   initialState,
			decisionHandler: decisionHandler,
		},
	}, nil
}
