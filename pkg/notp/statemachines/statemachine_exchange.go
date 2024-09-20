
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
	"bytes"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
)

// ExchangePacket encapsulates the data structure for an exchange packet used in the protocol.
type ExchangePacket struct {
	notppackets.Packet
}

// PacketType returns the specific type identifier for the exchange packet.
func (p *ExchangePacket) PacketType() int32 {
	return ExchangePacketType
}

// Marshal converts the ExchangePacket into a serialized byte slice for transmission.
func (p *ExchangePacket) Marshal() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	return buffer.Bytes(), nil
}

// Unmarshal populates the ExchangePacket with data from the given byte slice.
func (p *ExchangePacket) Unmarshal(data []byte) error {
	// buffer := bytes.NewBuffer(data)
	return nil
}
