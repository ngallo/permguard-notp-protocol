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
func createStatePacket(flow FlowType, currentStateID uint16, messageCode uint16, messageValue uint64) (*notpsmpackets.StatePacket, *HandlerContext, error) {
	handlerCtx := &HandlerContext{
		flow:           flow,
		currentStateID: currentStateID,
	}
	packet := &notpsmpackets.StatePacket{
		MessageCode:  messageCode,
		MessageValue: messageValue,
		ErrorCode:    0,
	}
	return packet, handlerCtx, nil
}

// createAndHandleStatePacket creates a state packet and handles it.
func createAndHandleStatePacket(runtime *StateMachineRuntimeContext, messageCode uint16, messageValue uint64, packetables []notppackets.Packetable) (*notpsmpackets.StatePacket, []notppackets.Packetable, error) {
	statePacket, handlerCtx, err := createStatePacket(runtime.GetFlowType(), runtime.currentStateID, messageCode, messageValue)
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to create state packet: %w", err)
	}
	var handledPacketables []notppackets.Packetable
	if statePacket.MessageCode != notpsmpackets.ActionResponseMessage {
		handlerReturn, err := runtime.HandleStream(handlerCtx, statePacket, packetables)
		handledPacketables = handlerReturn.Packetables
		if err != nil {
			return nil, nil, fmt.Errorf("notp: failed to handle created packet: %w", err)
		}
		statePacket.MessageValue = handlerReturn.MessageValue
		statePacket.ErrorCode = handlerReturn.ErrorCode
	}
	return statePacket, handledPacketables, nil
}

// createAndHandleAndStreamStatePacket creates a state packet, handles it, and streams it.
func createAndHandleAndStreamStatePacket(runtime *StateMachineRuntimeContext, messageCode uint16, packetables []notppackets.Packetable) error {
	messageValue := notppackets.CombineUint32toUint64(notpsmpackets.UnknownValue, notpsmpackets.UnknownValue)
	return createAndHandleAndStreamStatePacketWithValue(runtime, messageCode, messageValue, packetables)
}

// createAndHandleAndStreamStatePacketWithValue creates a state packet with value, handles it, and streams it.
func createAndHandleAndStreamStatePacketWithValue(runtime *StateMachineRuntimeContext, messageCode uint16, messageValue uint64, packetables []notppackets.Packetable) error {
	packet, packetables, err := createAndHandleStatePacket(runtime, messageCode, messageValue, packetables)
	if err != nil {
		return fmt.Errorf("notp: failed to create and handle packet: %w", err)
	}
	streamPacketables := append([]notppackets.Packetable{packet}, packetables...)
	runtime.SendStream(streamPacketables)
	return nil
}

// receiveAndHandleStatePacket receives a state packet and handles it.
func receiveAndHandleStatePacket(runtime *StateMachineRuntimeContext, expectedState uint16) (*notpsmpackets.StatePacket, []notppackets.Packetable, error) {
	handlerCtx := &HandlerContext{
		flow:           runtime.GetFlowType(),
		currentStateID: runtime.GetCurrentStateID(),
	}
	packetsStream, err := runtime.ReceiveStream()
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to receive packets: %w", err)
	}
	statePacket := &notpsmpackets.StatePacket{}
	data, err := packetsStream[0].Serialize()
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to serialize packet: %w", err)
	}
	err = statePacket.Deserialize(data)
	if err != nil {
		return nil, nil, fmt.Errorf("notp: failed to deserialize state packet: %w", err)
	}
	if statePacket.HasError() {
		return nil, nil, fmt.Errorf("notp: received state packet with error: %d", statePacket.ErrorCode)
	}
	if statePacket.MessageCode != expectedState {
		return nil, nil, fmt.Errorf("notp: received unexpected state code: %d", statePacket.MessageCode)
	}
	var handledPacketables []notppackets.Packetable
	if statePacket.MessageCode != notpsmpackets.ActionResponseMessage {
		handlerReturn, err := runtime.HandleStream(handlerCtx, statePacket, packetsStream[1:])
		handledPacketables = handlerReturn.Packetables
		if err != nil {
			return nil, nil, fmt.Errorf("notp: failed to handle created packet: %w", err)
		}
		statePacket.MessageValue = handlerReturn.MessageValue
		statePacket.ErrorCode = handlerReturn.ErrorCode
	}
	return statePacket, handledPacketables, nil
}
