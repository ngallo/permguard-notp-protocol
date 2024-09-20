
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

// OperationType represents the type of operation that the NOTP protocol is performing.
type OperationType string

const (
	// AdvertisementPacketType represents the type of the advertisement packet.
	AdvertisementPacketType = uint32(1)
	// NegotiationPacketType represents the type of the negotiation packet.
	NegotiationPacketType = uint32(2)
	// ExchangePacketType represents the type of the exchange packet.
	ExchangePacketType 	= uint32(3)

	// PushOperation represents the push operation type.
    PushOperation    OperationType = "push"
	// PullOperation represents the pull operation type.
    PullOperation    OperationType = "pull"
	// DefaultOperation represents the default operation type.
    DefaultOperation OperationType = PushOperation
)
