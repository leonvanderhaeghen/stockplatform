package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productv1 "github.com/leonvanderhaeghen/stockplatform/gen/go/product/v1"
)

func main() {
	// Set up a connection to the server
	addr := "127.0.0.1:50053"
	log.Printf("Connecting to gRPC server at %s...\n", addr)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to gRPC server")

	// Create a client
	client := productv1.NewProductServiceClient(conn)

	// Test GetProduct with a test ID
	testGetProduct(client, "test-product-id")
}

func testGetProduct(client productv1.ProductServiceClient, id string) {
	log.Printf("Testing GetProduct with ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProduct(ctx, &productv1.GetProductRequest{Id: id})
	if err != nil {
		log.Printf("GetProduct failed: %v\n", err)
		return
	}

	// Pretty print the response
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal response: %v\n", err)
		return
	}

	log.Printf("GetProduct response:\n%s\n", string(jsonBytes))
}
