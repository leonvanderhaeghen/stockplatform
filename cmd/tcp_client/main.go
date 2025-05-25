package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Server address
	addr := "localhost:50053"

	// Try to establish a TCP connection
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", addr, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Successfully connected to %s\n", addr)
	fmt.Println("TCP connection test passed!")
}
