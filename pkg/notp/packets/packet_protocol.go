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
)

// ProtocolPacket represents a protocol packet.
type ProtocolPacket struct {
	Version int32
}

// GetType returns the type of the packet.
func (p *ProtocolPacket) GetType() int32 {
	return 0
}

// Serialize serializes the packet.
func (p *ProtocolPacket) Serialize() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.LittleEndian, p.Version); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Deserialize deserializes the packet.
func (p *ProtocolPacket) Deserialize(data []byte) error {
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.LittleEndian, &p.Version); err != nil {
		return err
	}
	return nil
}
