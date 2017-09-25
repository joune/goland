package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joune/zenly/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address            = "localhost:50051"
	defaultUser uint64 = 1
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := data.NewCategoriesClient(conn)

	// Get user ID from args
	user := defaultUser
	if len(os.Args) > 1 {
		if user, err = strconv.ParseUint(os.Args[1], 10, 64); err != nil {
			panic(err)
		}
	}
	// Contact the server and print out its response.
	resp, err := c.GetCategories(context.Background(), &data.CategoriesRequest{User: user})
	if err != nil {
		log.Fatalf("failed to fetch friend categories: %v", err)
	}
	log.Printf("User%v Categories: %v", user, resp)
}
