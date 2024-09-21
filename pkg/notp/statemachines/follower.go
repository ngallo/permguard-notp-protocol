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

// followerPullAdvertiseState handles the pull advertisement phase in the protocol.
func followerPullAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	_, advPacket, packetables, err := CrateAndHandlePacket(runtime, PullStateMachineType, false, notpsmpackets.ClientAdvertiseRequestChanges, notpsmpackets.AlgoFetchAll,
		func(basePacket *notpsmpackets.BasePacket) notppackets.Packetable {
			return &notpsmpackets.AdvertisementPacket{BasePacket: *basePacket}
	})
	if err != nil {
		return false, nil, fmt.Errorf("notp: failed to create and handle advertisement packet: %w", err)
	}
	runtime.SendStream(append([]notppackets.Packetable{advPacket}, packetables...))
	return false, FollowerNegotiateState, nil
}

// followerPushAdvertiseState handles the push advertisement phase in the protocol.
func followerPushAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, FollowerNegotiateState, nil
}

// FollowerAdvertiseState handles the advertisement phase in the protocol.
func FollowerAdvertiseState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	if runtime.GetOperation() == PullStateMachineType {
		return followerPullAdvertiseState(runtime)
	} else if runtime.GetOperation() == PushStateMachineType {
		return followerPushAdvertiseState(runtime)
	}
	return false, nil, fmt.Errorf("notp: invalid operation type: %s", runtime.GetOperation())
}

// FollowerNegotiateState manages the negotiation phase in the protocol.
func FollowerNegotiateState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, FollowerExchangeState, nil
}

// FollowerExchangeState governs the exchange phase in the protocol.
func FollowerExchangeState(runtime *StateMachineRuntimeContext) (bool, StateTransitionFunc, error) {
	return false, FinalState, nil
}
