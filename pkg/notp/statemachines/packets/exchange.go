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

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
)

// ExchangePacket encapsulates the data structure for an exchange packet used in the protocol.
type ExchangePacket struct {
}

// GetType returns the packet type.
func (p *ExchangePacket) GetType() uint64 {
	return notppackets.CombineUint32toUint64(ExchangePacketType, 0)
}

// Serialize serializes the packet.
func (p *ExchangePacket) Serialize() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	return buffer.Bytes(), nil
}

// Deserialize deserializes the packet.
func (p *ExchangePacket) Deserialize(data []byte) error {
	// buffer := bytes.NewBuffer(data)
	return nil
}
