package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestResolveUDPAddr(t *testing.T) {
	target := "1.1.1.1:80"
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		fmt.Println("Error resolving target:", err)
		return
	}

	fmt.Printf("Resolved target: IP: %s, Port: %d\n", addr.IP, addr.Port)
}

func TestUdpFloodSend(t *testing.T) {
	// use local host for test purposes
	//127.0.0.0/8 is reserved for loopback addresses
	receiverAddr := "127.0.0.1:8080"
	senderAddr := "127.0.0.1:8081"
	recvBuf := make([]byte, 1024)

	// Set up a local UDP receiver server to receive packets
	receiver, err := net.ListenPacket("udp", receiverAddr)
	if err != nil {
		t.Fatalf("Error setting up UDP receiver server: %v", err)
	}
	defer receiver.Close()

	// Set up a local UDP sender server to send packets
	sender, err := net.ListenPacket("udp", senderAddr)
	if err != nil {
		t.Fatalf("Error setting up UDP sender server: %v", err)
	}
	defer sender.Close()

	// Resolve the receiver address for the sender to send packets
	receiverUDPaddr, err := net.ResolveUDPAddr("udp", receiverAddr)
	if err != nil {
		t.Fatalf("Error resolving receiver address: %v", err)
	}

	// Set up the test payload and number of packets to send
	payload := randomPayload(1024)
	numpackets := 1

	// Run the udpFlood function using the sender server
	go udpflood_for_testing(sender, receiverUDPaddr, numpackets, payload)

	// Read packets received by the UDP receiver server
	n, _, err := receiver.ReadFrom(recvBuf)
	if err != nil {
		t.Fatalf("Error reading from receiver server: %v", err)
	}

	// Print the received UDP packet
	fmt.Printf("Received UDP packet: %s\n", recvBuf[:n])

	// Check if the received payload matches the sent payload
	if !bytes.Equal(recvBuf[:n], payload) {
		t.Errorf("Received payload does not match sent payload")
	}
}
