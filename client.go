package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type PacketData struct {
	Packet Packet
	Data   []byte
	Addr   net.UDPAddr
}

func RunClient(server_ip string) {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	data := ReconcilliationData{"Hello, server!"}

	packet := Packet{}
	packet.PacketType = PacketTypeMatchFind

	other_addr := net.UDPAddr{IP: net.ParseIP(server_ip), Port: 8080}

	raw_data, _ := SerializePacket(packet, data)
	_, err = conn.WriteToUDP(raw_data, &other_addr)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Data sent")
	packet_channel := make(chan PacketData)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error reading", err)
			}

			packet, data, _ := DeserializePacket(buf[:n])
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
			case PacketTypeMatchConnect:
				err = dec.Decode(&other_addr)

				fmt.Println(packet, other_addr)

				packet = Packet{}
				packet.PacketType = PacketTypeNegotiate
				data = ReconcilliationData{"Hey other client!"}

				raw_data, err := SerializePacket(packet, data)
				if err != nil {
					fmt.Println("error serializing packet", err)
				}

				_, err = conn.WriteToUDP(raw_data, &other_addr)
				if err != nil {
					fmt.Println("something went wrong when reaching out to match", err)
				}
			case PacketTypeNegotiate:
				var inner_data ReconcilliationData
				err = dec.Decode(&inner_data)

				// if we get this packet there is a presumption that we have already
				// broken through the NAT address by sending a packet to said address.

				// therefore we can safely assume that the incomming packet is from the owner we want to connect with
				// and then we can set the owner of the packet to our desired target address to assert the case
				other_addr = packet_data.Addr

				fmt.Println(packet_data.Packet, inner_data)
			}

		case <-time.After(5 * time.Second):
			fmt.Println("sending keepalive packet to", other_addr)

			packet = Packet{}
			packet.PacketType = PacketTypeKeepAlive
			data = ReconcilliationData{"keepalive"}

			serialized_packet, _ := SerializePacket(packet, data)

			_, err = conn.WriteToUDP(serialized_packet , &other_addr)
			if err != nil {
				fmt.Println("something went wrong when reaching out to match", err)
			}
		}
	}
}
