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
	"errors"
	"unsafe"
)

const (
	// PacketNullByte is the null byte used to separate data in the packet.
	PacketNullByte = 0x01
)

// writeStreamDataPacket writes a stream data packet to the buffer.
func writeStreamDataPacket(data []byte, packetType *int32, packetStream *int32, payload []byte) ([]byte, error) {
	size := int32(len(payload))
	if packetType != nil {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, *packetType); err != nil {
			return nil, err
		}
		data = append(data, buf.Bytes()...)
	}
	if packetStream != nil {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, *packetStream); err != nil {
			return nil, err
		}
		data = append(data, buf.Bytes()...)
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, size); err != nil {
		return nil, err
	}
	data = append(data, buf.Bytes()...)
	data = append(data, PacketNullByte)
	data = append(data, payload...)
	return data, nil
}

// writeDataPacket writes a data packet to the buffer.
func writeDataPacket(data []byte, payload []byte) ([]byte, error) {
	return writeStreamDataPacket(data, nil, nil, payload)
}

// indexDataStreamPacket indexes a stream data packet in the buffer.
func indexDataStreamPacket(offset int, data []byte) (int, int, int32, int32, error) {
	data = data[offset:]
	delimiterIndex := bytes.IndexByte(data, PacketNullByte)
	if delimiterIndex == -1 {
		return -1, -1, -1, -1, errors.New("notp: delimiter not found")
	}
	headerData := data[:delimiterIndex]
	idSize := int(unsafe.Sizeof(int32(0)))
	if len(headerData) != idSize * 3 {
		return -1, -1, -1, -1, errors.New("notp: invalid data: missing or invalid header")
	}
	dataOffset := delimiterIndex + 1
	values := []int32{0, 0, 0}
	for count := range values {
		start := idSize * count
		end := (idSize * count) + idSize
		values[count] = int32(binary.LittleEndian.Uint32(headerData[start:end]))
	}
	packetType := values[0]
	packetStream := values[1]
	size := int(values[2])
	return offset + dataOffset, size, packetType, packetStream, nil
}

// readStreamDataPacket reads a stream data packet from the buffer.
func readStreamDataPacket(offset int, data []byte) ([]byte, int, int, int32, int32, error) {
	offset, size, packetType, packetStream, err := indexDataStreamPacket(offset, data)
	if err != nil {
		return nil, -1, -1, -1, -1, err
	}
	payload := data[offset:offset + size]
	return payload, offset, size, packetType, packetStream, nil
}

// indexDataPacket indexes a data packet in the buffer.
func indexDataPacket(offset int, data []byte) (int, int, error) {
	data = data[offset:]
	delimiterIndex := bytes.IndexByte(data, PacketNullByte)
	if delimiterIndex == -1 {
		return -1, -1, errors.New("notp: delimiter not found")
	}
	headerData := data[:delimiterIndex]
	idSize := int(unsafe.Sizeof(int32(0)))
	if len(headerData) != idSize {
		return -1, -1, errors.New("notp: invalid data: missing or invalid header")
	}
	dataOffset := delimiterIndex + 1
	size := int(binary.LittleEndian.Uint32(headerData))
	return offset + dataOffset, size, nil
}

// readDataPacket reads a data packet from the buffer.
func readDataPacket(offset int, data []byte) ([]byte, int, int, error) {
	offset, size, err := indexDataPacket(offset, data)
	if err != nil {
		return nil, -1, -1, err
	}
	payload := data[offset:offset + size]
	return payload, offset, size, nil
}

