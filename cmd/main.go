package main

import (
	"fmt"
	"log"
	"net"

	"github.com/songgao/packets/ethernet"
)

func main() {
	// Initialize MAC table
	var MacTable = make(map[string]*net.UDPAddr)

	// Create UDP address
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Create UDP connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("UDP server listening on :8080")

	for {
		buffer := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		// Parse ethernet frame
		frame := ethernet.Frame(buffer[:n])
		sourceMac := frame.Source().String()
		desMac := frame.Destination().String()

		// Update MAC table with source address
		MacTable[sourceMac] = remoteAddr
		log.Printf("Learned MAC address %s from %s", sourceMac, remoteAddr.String())

		if desMac == "ff:ff:ff:ff:ff:ff" {
			// Broadcast packet
			log.Println("Broadcast packet")

			for _, addr := range MacTable {
				conn.WriteToUDP(buffer[:n], addr)
			}
		} else {
			fmt.Println("Sending packet to", MacTable[desMac].String())
			conn.WriteToUDP(buffer[:n], MacTable[desMac])
		}

	}
}
