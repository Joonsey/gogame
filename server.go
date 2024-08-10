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

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	fmt.Println("Listening")
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		packet, data, err := DeserializePacket(buf[:n])
		if packet.PacketType == 1 {
			var inner_data InnerData
			dec := gob.NewDecoder(bytes.NewReader(data))
			err := dec.Decode(&inner_data)
			if err != nil {
				fmt.Println("error during decoding", err)
			}

			fmt.Println(packet, inner_data)
		}
	}
}
