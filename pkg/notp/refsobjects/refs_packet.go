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

package refsobjects

import (
	"bytes"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
)

// RefsObjPacket represents a reference object packet.
type RefsObjPacket struct {
	notppackets.Packet
}

// GetType returns the type of the packet.
func (p *RefsObjPacket) GetType() int32 {
	return 0
}

// Serialize serializes the packet.
func (p *RefsObjPacket) Serialize() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	return buffer.Bytes(), nil
}

// Deserialize deserializes the packet.
func (p *RefsObjPacket) Deserialize(data []byte) error {
	// buffer := bytes.NewBuffer(data)
	return nil
}
