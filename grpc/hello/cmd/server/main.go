package main

import (
	"context"
	"log"
	"net"

	idl "github.com/dipperkun/go-notes/grpc/hello/idl"
	"google.golang.org/grpc"
)

// server is used to implement idl.GreeterServer.
type server struct {
	idl.UnimplementedGreeterServer
}

// One implements idl.GreeterServer
func (s *server) One(ctx context.Context, in *idl.Req) (*idl.Resp, error) {
	log.Printf("Received: %v", in.GetName())
	return &idl.Resp{Gift: "Hello " + in.GetName()}, nil
}

func main() {
	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	idl.RegisterGreeterServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
