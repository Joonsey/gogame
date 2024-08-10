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

func timeoutStaleConnections(keyword_map* map[string]Connection) {
	for key, value := range map[string]Connection(*keyword_map) {
		if time.Now().UnixMilli() - value.Time > 7000 {
			fmt.Printf("%s user timed out using '%s' connection key\n", value.Addr, value.Keyword)
			delete(*keyword_map, key)
		}
	}
}


func RunServer() {
	server_addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	var keyword_map map[string]Connection
	keyword_map = make(map[string]Connection)

	conn, err := net.ListenUDP("udp", server_addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	fmt.Println("Listening")
	defer conn.Close()

	packet_channel := make(chan PacketData)

	go func() {
		for {
			timeoutStaleConnections(&keyword_map)
			time.Sleep(time.Second * 1)
		}
	}()


	go func() {
		buf := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error reading", err)
			}

			packet, data, err := DeserializePacket(buf[:n])
			if err != nil {
				fmt.Println("error reading", err)
			}

			packet_data := PacketData{packet, data, *addr}
			packet_channel <- packet_data
		}
	}()

	for {
		select {
		case packet_data := <-packet_channel:
			dec := gob.NewDecoder(bytes.NewReader(packet_data.Data))
			switch packet_data.Packet.PacketType {
			case PacketTypeKeepAlive:
				for key, value := range keyword_map {
					if value.Addr.String() == packet_data.Addr.String() {
						value.Time = time.Now().UnixMilli()
						keyword_map[key] = value
					}
				}

			case PacketTypeMatchFind:
				var inner_data ReconcilliationData
				err := dec.Decode(&inner_data)
				if err != nil {
					fmt.Println("error during decoding", err)
				}

				if keyword_map[inner_data.Name].Keyword != "" {
					packet := Packet{}
					packet.PacketType = PacketTypeMatchConnect
					data := packet_data.Addr
					serialized_packet, err := SerializePacket(packet, data)
					if err != nil {
						fmt.Println("error during serialization", err)
					}

					conn.WriteToUDP(serialized_packet, keyword_map[inner_data.Name].Addr)

					data = *keyword_map[inner_data.Name].Addr
					serialized_packet, err = SerializePacket(packet, data)
					if err != nil {
						fmt.Println("error during serialization", err)
					}
					conn.WriteToUDP(serialized_packet, &packet_data.Addr)

					delete(keyword_map, inner_data.Name)

				} else {
					keyword_map[inner_data.Name] = Connection{inner_data.Name, &packet_data.Addr, time.Now().UnixMilli()}
				}

				fmt.Println(packet_data.Packet, packet_data.Addr, inner_data)
			}
		}
	}
}
