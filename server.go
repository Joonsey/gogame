package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

func RunServer() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	var keyword_map map[string]*net.UDPAddr
	keyword_map = make(map[string]*net.UDPAddr)

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	fmt.Println("Listening")
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		packet, data, err := DeserializePacket(buf[:n])
		if packet.PacketType == PacketTypeMatchFind {
			var inner_data InnerData
			dec := gob.NewDecoder(bytes.NewReader(data))
			err := dec.Decode(&inner_data)
			if err != nil {
				fmt.Println("error during decoding", err)
			}

			if keyword_map[inner_data.Name] != nil {
				packet := Packet{}
				packet.PacketType = PacketTypeMatchConnect
				data := *addr
				packet_data, err := SerializePacket(packet, data)
				if err != nil {
					fmt.Println("error during serialization", err)
				}
				conn.WriteToUDP(packet_data, keyword_map[inner_data.Name])

				data = *keyword_map[inner_data.Name]
				packet_data, err = SerializePacket(packet, data)
				if err != nil {
					fmt.Println("error during serialization", err)
				}
				conn.WriteToUDP(packet_data, addr)

				delete(keyword_map, inner_data.Name)

			} else {
				keyword_map[inner_data.Name] = addr
			}

			fmt.Println(packet, addr, inner_data)
		}
	}
}
