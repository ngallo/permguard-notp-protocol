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

// Package transport implements the transport layer of the NOTP protocol.
package transport

import (
	"errors"

	notppackets "github.com/permguard/permguard-notp-protocol/pkg/notp/packets"
)

// InMemoryStream simulates an in-memory stream for packet transmission.
type InMemoryStream struct {
	packets []notppackets.Packet
}

// TransmitPacket appends a packet to the in-memory stream.
func (t *InMemoryStream) TransmitPacket(packet *notppackets.Packet) error {
	if packet == nil {
		return errors.New("notp: cannot transmit a nil packet")
	}
	t.packets = append(t.packets, *packet)
	return nil
}

// ReceivePacket retrieves the oldest packet from the in-memory stream.
func (t *InMemoryStream) ReceivePacket() (*notppackets.Packet, error) {
	if len(t.packets) == 0 {
		return nil, errors.New("notp: no packets available in transport layer")
	}
	packet := t.packets[0]
	t.packets = t.packets[1:]
	return &packet, nil
}

// NewInMemoryStream creates and initializes a new in-memory stream.
func NewInMemoryStream() (*InMemoryStream, error) {
	return &InMemoryStream{
		packets: make([]notppackets.Packet, 0),
	}, nil
}
