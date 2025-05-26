package main

import (
	"log"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
)

func main() {
	// Set up a connection to the server
	addr := "localhost:50053"
	timeout := 5 * time.Second
	
	conn, client, err := grpcclient.NewProductClient(addr, timeout)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	// Test GetProduct with a test ID
	_, err = grpcclient.TestGetProduct(client, "test-product-id")
	if err != nil {
		log.Fatalf("GetProduct failed: %v", err)
	}
}
