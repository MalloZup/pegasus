package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/MalloZup/pegasus/helloworld"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultName     = "world"
	retriesParallel = 100
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "second"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "third"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "quattro"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "cinque"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "sei"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "sette"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "otto"})
	go c.SayHelloAgain(ctx, &pb.HelloRequest{Name: "nove"})

	log.Printf("Greeting: %s", r.GetMessage())
}