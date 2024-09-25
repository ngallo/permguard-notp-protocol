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

	// AlgoFetchAll represents fetching all available data.
	AlgoFetchAll = uint16(101)
	// AlgoFetchExactVersion represents fetching an exact version of the data.
	AlgoFetchExactVersion = uint16(100)
	// AlgoFetchMinimal represents fetching minimal data when bandwidth is limited.
	AlgoFetchMinimal = uint16(1002)

	// Represents a request to obtain the current state.
	RequestCurrentState = uint16(111)
	// Respond with the current state.
	RespondCurrentState = uint16(112)
	// Notify other parties about the current state.
	NotifyCurrentState = uint16(113)

	// Submit a request to initiate a negotiation process.
	SubmitNegotiationRequest = uint16(114)
	// Respond to reject the submitted negotiation request.
	RejectNegotiationRequest = uint16(115)
	// Approve the submitted negotiation request.
	ApproveNegotiationRequest = uint16(116)

	// Start the data streaming process.
	InitiateDataStream = uint16(117)
	// Continue sending data within the ongoing stream.
	ContinueDataStream = uint16(118)
	// Conclude the data streaming process.
	ConcludeDataStream = uint16(119)
)

// StatePacket encapsulates the data structure for a base packet used in the protocol.
type StatePacket struct {
	StateCode     uint16
	AlgorithmCode uint16
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

	err = binary.Write(buffer, binary.BigEndian, p.AlgorithmCode)
	if err != nil {
		return nil, fmt.Errorf("failed to write AlgorithmCode: %v", err)
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

	err = binary.Read(buffer, binary.BigEndian, &p.AlgorithmCode)
	if err != nil {
		return fmt.Errorf("failed to read AlgorithmCode: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &p.ErrorCode)
	if err != nil {
		return fmt.Errorf("failed to read ErrorCode: %v", err)
	}

	return nil
}
