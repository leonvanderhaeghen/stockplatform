package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	productv1 "github.com/yourusername/stockplatform/gen/go/product/v1"
)

func main() {
	// Set up a connection to the server
	log.Println("Connecting to gRPC server at 127.0.0.1:50053...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	conn, err := grpc.DialContext(ctx, "127.0.0.1:50053",
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

	// Test GetProduct with a non-existent ID
	testGetProduct(client, "non-existent-id")
}

func testGetProduct(client productv1.ProductServiceClient, id string) {
	log.Printf("Testing GetProduct with ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProduct(ctx, &productv1.GetProductRequest{Id: id})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			log.Printf("gRPC error: %v (code: %s, message: %s)\n", 
				err, status.Code(), status.Message())
		} else {
			log.Printf("Non-gRPC error: %v\n", err)
		}
		return
	}

	log.Printf("GetProduct response: %+v\n", resp)
}
