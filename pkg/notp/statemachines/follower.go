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

	notptransport "github.com/permguard/permguard-notp-protocol/pkg/notp/transport"
)

// NewFollowerStateMachine creates and configures a new follower state machine for the given operation.
func NewFollowerStateMachine(operation OperationType, decisionHandler DecisionHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
    if operation == "" {
        operation = DefaultOperation
    }
    stateMachine, err := NewStateMachine(operation, FollowerAdvertiseState, decisionHandler, transportLayer)
    if err != nil {
        return nil, fmt.Errorf("notp: failed to create follower state machine: %w", err)
    }
    return stateMachine, nil
}

// FollowerAdvertiseState handles the advertisement phase in the protocol.
func FollowerAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, FollowerExchangeState, nil
}

// FollowerNegotiateState manages the negotiation phase in the protocol.
func FollowerNegotiateState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, FollowerExchangeState, nil
}

// FollowerExchangeState governs the exchange phase in the protocol.
func FollowerExchangeState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, FinalState, nil
}

