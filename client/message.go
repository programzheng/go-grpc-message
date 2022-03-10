package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/programzheng/go-grpc-message"
	"google.golang.org/grpc"
)

const (
	defaultMessage = "hello"
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	message = flag.String("message", defaultMessage, "Message")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayMessage(ctx, &pb.MessageRequest{Message: *message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
