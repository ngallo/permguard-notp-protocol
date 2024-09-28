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

package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
)

const (
	// StatePacketType represents the type of the state packet.
	StatePacketType = uint32(10)

	// NotifyCurrentState represents the notification of the current state.
	NotifyCurrentState = uint16(101)
	// RequestCurrentState represents the request for the current state.
	RequestCurrentState = uint16(102)
	// RespondCurrentState represents the response to the request for the current state.
	RespondCurrentState = uint16(103)

	// SubmitNegotiationRequest represents the submission of a negotiation request.
	SubmitNegotiationRequest = uint16(124)
	// RespondNegotiationRequest represents the response to a negotiation request.
	RespondNegotiationRequest = uint16(125)

	// ExchangeDataStream represents the exchange of data stream.
	ExchangeDataStream = uint16(146)
)

// StatePacket encapsulates the data structure for a base packet used in the protocol.
type StatePacket struct {
	StateCode     uint16
	ErrorCode     uint16
}

// GetType returns the packet type.
func (p *StatePacket) GetType() uint64 {
	return notppackets.CombineUint32toUint64(StatePacketType, 0)
}

// HasError returns true if the packet has errors.
func (p *StatePacket) HasError() bool {
	return p.ErrorCode != 0
}

// Serialize serializes the packet into bytes.
func (p *StatePacket) Serialize() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, p.StateCode)
	if err != nil {
		return nil, fmt.Errorf("failed to write OperationCode: %v", err)
	}

	err = binary.Write(buffer, binary.BigEndian, p.ErrorCode)
	if err != nil {
		return nil, fmt.Errorf("failed to write ErrorCode: %v", err)
	}

	return buffer.Bytes(), nil
}

// Deserialize deserializes the packet from bytes.
func (p *StatePacket) Deserialize(data []byte) error {
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &p.StateCode)
	if err != nil {
		return fmt.Errorf("failed to read OperationCode: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &p.ErrorCode)
	if err != nil {
		return fmt.Errorf("failed to read ErrorCode: %v", err)
	}

	return nil
}
