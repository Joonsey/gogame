package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Connection struct {
	Keyword string
	Addr *net.UDPAddr
	Time int64
}

func RunServer() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	var keyword_map map[string]Connection
	keyword_map = make(map[string]Connection)

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

		for key, value := range map[string]Connection(keyword_map) {
			if time.Now().UnixMilli() - value.Time > 7000 {
				fmt.Printf("%s user timed out using '%s' connection key\n", value.Addr, value.Keyword)
				delete(keyword_map, key)
			}
		}

		packet, data, err := DeserializePacket(buf[:n])
		if packet.PacketType == PacketTypeKeepAlive {
			for key, value := range keyword_map {
				if value.Addr.String() == addr.String() {
					fmt.Println(value)
					value.Time = time.Now().UnixMilli()
					keyword_map[key] = value
					fmt.Printf("%s user refreshed\n", value.Addr)
				}
			}
		}

		if packet.PacketType == PacketTypeMatchFind {
			var inner_data ReconcilliationData
			dec := gob.NewDecoder(bytes.NewReader(data))
			err := dec.Decode(&inner_data)
			if err != nil {
				fmt.Println("error during decoding", err)
			}

			if keyword_map[inner_data.Name].Keyword != "" {
				packet := Packet{}
				packet.PacketType = PacketTypeMatchConnect
				data := *addr
				packet_data, err := SerializePacket(packet, data)
				if err != nil {
					fmt.Println("error during serialization", err)
				}

				conn.WriteToUDP(packet_data, keyword_map[inner_data.Name].Addr)

				data = *keyword_map[inner_data.Name].Addr
				packet_data, err = SerializePacket(packet, data)
				if err != nil {
					fmt.Println("error during serialization", err)
				}
				conn.WriteToUDP(packet_data, addr)

				delete(keyword_map, inner_data.Name)

			} else {
				keyword_map[inner_data.Name] = Connection{inner_data.Name, addr, time.Now().UnixMilli()}
			}

			fmt.Println(packet, addr, inner_data)
		}
	}
}
