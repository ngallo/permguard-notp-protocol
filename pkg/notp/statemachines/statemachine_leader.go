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

// NewLeaderStateMachine initializes and returns a new leader state machine for the specified operation.
func NewLeaderStateMachine(operation OperationType, decisionHandler DecisionHandler, transportLayer *notptransport.TransportLayer) (*StateMachine, error) {
    var initialState StateTransitionFunc
    if operation == "" {
        operation = DefaultOperation
    }
    switch operation {
    case PushOperation:
        initialState = LeaderAdvertiseRequiredObjectsState
    case PullOperation:
        initialState = LeaderAdvertiseLatestObjectsState
    default:
        return nil, fmt.Errorf("notp: invalid operation: %s", operation)
    }

    stateMachine, err := NewStateMachine(initialState, decisionHandler, transportLayer)
    if err != nil {
        return nil, fmt.Errorf("notp: failed to create leader state machine: %w", err)
    }
    return stateMachine, nil
}

// LeaderAdvertiseRequiredObjectsState advertises the required objects to the leader.
func LeaderAdvertiseRequiredObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, LeaderNegotiatingState, nil
}

// LeaderAdvertiseLatestObjectsState advertises the latest objects to the leader.
func LeaderAdvertiseLatestObjectsState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, LeaderNegotiatingState, nil
}

// LeaderNegotiatingState manages the negotiation phase between the leader and the leader.
func LeaderNegotiatingState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, LeaderObjectExchangeState, nil
}

// LeaderObjectExchangeState manages the object exchange phase between the leader and the leader.
func LeaderObjectExchangeState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
    return false, FinalState, nil
}
