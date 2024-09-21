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
)

const (
	// AdvertisementPacketType represents the type of the advertisement packet.
	AdvertisementPacketType = uint32(10)
	// NegotiationPacketType represents the type of the negotiation packet.
	NegotiationPacketType = uint32(11)
	// ExchangePacketType represents the type of the exchange packet.
	ExchangePacketType 	= uint32(12)

	// AlgoFetchExactVersion represents fetching an exact version of the data.
	AlgoFetchExactVersion = uint16(100)
	// AlgoFetchAll represents fetching all available data.
	AlgoFetchAll = uint16(101)
	// AlgoFetchMinimal represents fetching minimal data when bandwidth is limited.
	AlgoFetchMinimal = uint16(1002)

	// ClientAdvertiseRequestLatestState represents the request latest state operation.
    ClientAdvertiseRequestLatestState = uint16(110)
	// ClientNegotiateRequestChangeset represents the request changeset operation.
    ClientNegotiateRequestChangeset = uint16(111)

	// ServerAdvertiseReplyLatestState represents the reply latest state operation.
	ServerAdvertiseReplyLatestState = uint16(210)
	// ServerNegotiateRespondToRequest represents the respond to request operation.
    ServerNegotiateRespondToRequest = uint16(211)
	// ServerExchangeStream represents the exchange stream operation.
    ServerExchangeStream = uint16(212)
)

// BasePacket encapsulates the data structure for a base packet used in the protocol.
type BasePacket struct {
    OperationCode uint16
    AlgorithmCode uint16
    ErrorCode     uint16
}

// HasError returns true if the packet has errors.
func (p *BasePacket) HasError() bool {
    return p.ErrorCode != 0
}

// Serialize serializes the packet into bytes.
func (p *BasePacket) Serialize() ([]byte, error) {
    buffer := bytes.NewBuffer([]byte{})

    // Serialize OperationCode (2 bytes, uint16)
    err := binary.Write(buffer, binary.BigEndian, p.OperationCode)
    if err != nil {
        return nil, fmt.Errorf("failed to write OperationCode: %v", err)
    }

    // Serialize AlgorithmCode (2 bytes, uint16)
    err = binary.Write(buffer, binary.BigEndian, p.AlgorithmCode)
    if err != nil {
        return nil, fmt.Errorf("failed to write AlgorithmCode: %v", err)
    }

    // Serialize ErrorCode (2 bytes, uint16)
    err = binary.Write(buffer, binary.BigEndian, p.ErrorCode)
    if err != nil {
        return nil, fmt.Errorf("failed to write ErrorCode: %v", err)
    }

    return buffer.Bytes(), nil
}

// Deserialize deserializes the packet from bytes.
func (p *BasePacket) Deserialize(data []byte) error {
    buffer := bytes.NewBuffer(data)

    err := binary.Read(buffer, binary.BigEndian, &p.OperationCode)
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
