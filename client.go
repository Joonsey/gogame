package main

import (
	"fmt"
	"net"
)

func RunClient() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	})
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	data := InnerData{"Hello, server!", 4}

	packet := Packet{}
	packet.PacketType = 1

	rawData, _ := SerializePacket(packet, data)
	_, err = conn.Write(rawData)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Data sent")
}
