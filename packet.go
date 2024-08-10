package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
)

type Packet struct {
	PacketType uint8
	HeaderSize uint32
	MagicBytes uint32
	PayloadSize uint32
	TotalSize  uint32
}

const MAGICBYTES = 73458339

type InnerData struct {
	Name string
	Id int
}

func ValidatePacket(packet Packet) error {
	if packet.TotalSize != packet.HeaderSize + packet.PayloadSize {
		return errors.New("packet has invalid sizes")
	}

	if packet.MagicBytes != MAGICBYTES {
		return errors.New("packet has invalid magic bytes")
	}

	return nil
}

func DeserializePacket(data []byte) (Packet, []byte, error) {
	var packet Packet
	r := bytes.NewReader(data)

	err := binary.Read(r, binary.BigEndian, &packet.PacketType)
	if err != nil {
		fmt.Println("error during decoding of packet type", err)
		return packet, nil, err
	}

	err = binary.Read(r, binary.BigEndian, &packet.HeaderSize)
	if err != nil {
		fmt.Println("error during decoding of header size", err)
		return packet, nil, err
	}

	err = binary.Read(r, binary.BigEndian, &packet.MagicBytes)
	if err != nil {
		fmt.Println("error during decoding of magic bytes", err)
		return packet, nil, err
	}

	err = binary.Read(r, binary.BigEndian, &packet.PayloadSize)
	if err != nil {
		fmt.Println("error during decoding of paylaod size", err)
		return packet, nil, err
	}

	err = binary.Read(r, binary.BigEndian, &packet.TotalSize)
	if err != nil {
		fmt.Println("error during decoding total size", err)
		return packet, nil, err
	}

	err = ValidatePacket(packet)
	if err != nil {
		fmt.Println("error during packet validation", err)
		return packet, nil, err
	}

	rawData := data[packet.HeaderSize:packet.TotalSize]
	return packet, rawData, nil
}

func SerializePacket(packet Packet, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	// setting metadata
	packet.HeaderSize = 17
	packet.MagicBytes = MAGICBYTES

	binary.Write(&buf, binary.BigEndian, packet.PacketType)
	binary.Write(&buf, binary.BigEndian, packet.HeaderSize)
	binary.Write(&buf, binary.BigEndian, packet.MagicBytes)

	dataBytes, err := SerializeData(data)
	if err != nil {
		return nil, err
	}
	packet.PayloadSize = uint32(len(dataBytes))
	packet.TotalSize = uint32(buf.Len() + 8) + uint32(len(dataBytes))
	// adding the 8 bytes from totalsize and payloadsize values

	binary.Write(&buf, binary.BigEndian, packet.PayloadSize)
	binary.Write(&buf, binary.BigEndian, packet.TotalSize)

	// Append encoded data
	buf.Write(dataBytes)

	return buf.Bytes(), nil
}

func SerializeData(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
