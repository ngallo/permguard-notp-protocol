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

// createStatePacket creates a state packet.
func createStatePacket(smType StateMachineType, state uint16) (*notpsmpackets.StatePacket, *HandlerContext, error) {
	handlerCtx := &HandlerContext{
		stateMachineType: smType,
	}
	packet := &notpsmpackets.StatePacket{
		StateCode: state,
		ErrorCode: 0,
	}
	return packet, handlerCtx, nil
}

// createAndHandleStatePacket creates a state packet and handles it.
func createAndHandleStatePacket(runtime *StateMachineRuntimeContext, state uint16, packetables []notppackets.Packetable) (bool, *notpsmpackets.StatePacket, []notppackets.Packetable, error) {
	packet, handlerCtx, err := createStatePacket(runtime.GetStateMachineType(), state)
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to create state packet: %w", err)
	}
	retry, handledPacketables, err := runtime.HandleStream(handlerCtx, packet, packetables)
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to handle created packet: %w", err)
	}
	return retry, packet, handledPacketables, nil
}

// createAndHandleAndStreamStatePacket creates a state packet and handles it.
func createAndHandleAndStreamStatePacket(runtime *StateMachineRuntimeContext, state uint16, packetables []notppackets.Packetable) error {
	_, packet, packetables, err := createAndHandleStatePacket(runtime, state, packetables)
	if err != nil {
		return fmt.Errorf("notp: failed to create and handle packet: %w", err)
	}
	streamPacketables := append([]notppackets.Packetable{packet}, packetables...)
	runtime.SendStream(streamPacketables)
	return nil
}

// receiveAndHandleStatePacket receives a state packet and handles it.
func receiveAndHandleStatePacket(runtime *StateMachineRuntimeContext, expectedState uint16) (bool, *notpsmpackets.StatePacket, []notppackets.Packetable, error) {
	handlerCtx := &HandlerContext{
		stateMachineType: runtime.GetStateMachineType(),
	}
	packetsStream, err := runtime.ReceiveStream()
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to receive packets: %w", err)
	}
	statePacket := &notpsmpackets.StatePacket{}
	data, err := packetsStream[0].Serialize()
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to serialize packet: %w", err)
	}
	err = statePacket.Deserialize(data)
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to deserialize state packet: %w", err)
	}
	if statePacket.HasError() {
		return false, nil, nil, fmt.Errorf("notp: received state packet with error: %d", statePacket.ErrorCode)
	}
	if statePacket.StateCode != expectedState {
		return false, nil, nil, fmt.Errorf("notp: received unexpected state code: %d", statePacket.StateCode)
	}
	retry, handledPacketables, err := runtime.Handle(handlerCtx, statePacket)
	if err != nil {
		return false, nil, nil, fmt.Errorf("notp: failed to handle created packet: %w", err)
	}
	return retry, statePacket, handledPacketables, nil
}
