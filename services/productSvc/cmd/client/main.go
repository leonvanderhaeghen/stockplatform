package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
)

func main() {
	serverAddr := flag.String("server", "localhost:50053", "The server address in the format of host:port")
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a new product service client
	client := productv1.NewProductServiceClient(conn)

	// Test ListProducts
	testListProducts(client)
}

func testListProducts(client productv1.ProductServiceClient) {
	fmt.Println("\n=== Testing ListProducts ===")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a filter
	filter := &productv1.ProductFilter{
		MinPrice:   10.0,  // Only products with price >= 10
		MaxPrice:   1000.0, // And price <= 1000
		SearchTerm: "",    // Optional: search term
	}

	// Create sort options
	sort := &productv1.ProductSort{
		Field: productv1.ProductSort_SORT_FIELD_PRICE,
		Order: productv1.ProductSort_SORT_ORDER_ASC,
	}

	// Create pagination
	pagination := &productv1.Pagination{
		Page:     1,
		PageSize: 10,
	}

	// Create the request
	req := &productv1.ListProductsRequest{
		Filter:     filter,
		Sort:       sort,
		Pagination: pagination,
	}

	// Call the ListProducts RPC
	resp, err := client.ListProducts(ctx, req)
	if err != nil {
		log.Fatalf("ListProducts failed: %v", err)
	}

	// Print the response
	fmt.Printf("Total products: %d\n", resp.TotalCount)
	fmt.Printf("Page: %d, Page Size: %d\n", resp.Page, resp.PageSize)
	fmt.Println("Products:")

	for i, product := range resp.Products {
		jsonData, _ := protojson.MarshalOptions{
			Multiline: true,
			Indent:    "  ",
		}.Marshal(product)
		fmt.Printf("%d. %s\n\n", i+1, string(jsonData))
	}
}
