// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	pb "github.com/MalloZup/pegasus/helloworld"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

// alwaysPassLimiter is an example limiter which implements Limiter interface.
// It does not limit any request because Limit function always returns false.
type rateLimiterInterceptor struct {
	MaxReq int
}

func (r *rateLimiterInterceptor) Limit() bool {
	fmt.Printf("%d:", r.MaxReq)
	if r.ReqCount >= r.MaxReq {
		return true
	}
	r.MaxReq++
	return false
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	limiter := &rateLimiterInterceptor{}
	limiter.ReqCount = 0
	limiter.MaxReq = 10
	s := grpc.NewServer(
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
