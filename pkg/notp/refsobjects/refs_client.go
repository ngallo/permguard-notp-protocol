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

package refsobjects

import (
	notpsmachine "github.com/permguard/permguard-notp-protocol/pkg/notp/statemachines"
	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// ClientStateMachine represents the client's state machine.
type ClientStateMachine struct {
	notpsmachine.StateMachine
}

// NewClientStateMachine creates and initializes a new state machine with the given initial state and transport layer.
func NewClientStateMachine(transportLayer *notptransport.TransportLayer) (*ClientStateMachine, error) {
	clientStateMachine := &ClientStateMachine{}
	stateMachine, err := notpsmachine.NewStateMachine(clientStateMachine.ClientAdvertisingState, transportLayer)
	if err != nil {
		return nil, err
	}
	clientStateMachine.StateMachine = *stateMachine
	return clientStateMachine, nil
}

// ClientAdvertisingState handles the client's actions during the advertising phase.
func (s *ClientStateMachine) ClientAdvertisingState(runtime *notpsmachine.StateMachineRuntimeContext) (bool, notpsmachine.StateTransitionFunc, error) {
	return false, s.ClientNegotiatingState, nil
}

// ClientNegotiatingState handles the client's actions during the negotiation phase.
func (s *ClientStateMachine) ClientNegotiatingState(runtime *notpsmachine.StateMachineRuntimeContext) (bool, notpsmachine.StateTransitionFunc, error) {
	return false, s.ClientObjectExchangeState, nil
}

// ClientObjectExchangeState handles the client's actions during the object exchange phase.
func (s *ClientStateMachine) ClientObjectExchangeState(runtime *notpsmachine.StateMachineRuntimeContext) (bool, notpsmachine.StateTransitionFunc, error) {
	return false, notpsmachine.FinalState, nil
}
