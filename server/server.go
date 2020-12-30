// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"

	"net"
	"time"

	pb "github.com/MalloZup/pegasus/helloworld"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/juju/ratelimit"
	"google.golang.org/grpc"
)

const (
	gatherTime    = 10 * time.Second
	port          = ":50051"
	tokenCapacity = 6
)

type rateLimiterInterceptor struct {
	TokenBucket *ratelimit.Bucket
}

func (r *rateLimiterInterceptor) Limit() bool {

	// if zero we reached rate limit, so return true ( report error to Grpc)
	tokenRes := r.TokenBucket.TakeAvailable(1)
	if tokenRes == 0 {
		fmt.Printf("Reached Rate-Limiting %d \n", r.TokenBucket.Available())
		return true
	}
	// debug
	fmt.Printf("Token Avail %d \n", r.TokenBucket.Available())

	// if tokenRes is not zero, means gRpc request can continue to flow without rate limiting :)
	return false
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer just some fake requests
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// some example request, no serious code here
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func main() {
	fmt.Println("started gRPC")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("listening...")
	limiter := &rateLimiterInterceptor{}

	limiter.TokenBucket = ratelimit.NewBucket(gatherTime, int64(tokenCapacity))
	s := grpc.NewServer(
		// init the Ratelimiting middleware
		grpc_middleware.WithUnaryServerChain(
			grpc_ratelimit.UnaryServerInterceptor(limiter),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ratelimit.StreamServerInterceptor(limiter),
		),
	)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
