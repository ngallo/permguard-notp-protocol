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

// Packet represents a packet.
type Packet struct {
	Data []byte
}

// PacketConverterHandler defines a function type for handling packet conversions.
type PacketConverterHandler func (packetType uint32, data []byte) (Packetable, error)

// Packetable represents a packet that can be serialized and deserialized.
type Packetable interface {
	GetType() uint64
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// CombineUint32toUint64 combina due uint32 in un singolo uint64.
func CombineUint32toUint64(high, low uint32) uint64 {
    return (uint64(high) << 32) | uint64(low)
}

// SplitUint64toUint32 suddivide un uint64 in due uint32.
func SplitUint64toUint32(value uint64) (uint32, uint32) {
    high := uint32(value >> 32)
    low := uint32(value & 0xFFFFFFFF)
    return high, low
}
