package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	idl "github.com/dipperkun/go-notes/grpc/hello/idl"
	"google.golang.org/grpc"
)

// server is used to implement idl.GreeterServer.
type server struct {
	idl.UnimplementedGreeterServer
}

// One implements idl.GreeterServer
// unary -> unary
func (s *server) One(ctx context.Context, in *idl.Req) (*idl.Resp, error) {
	log.Printf("Received: %v", in.GetName())
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	return &idl.Resp{Gift: "Hello " + in.GetName()}, nil
}

// Two implements idl.GreetServer
// unary -> stream
func (s *server) Two(in *idl.Req, strm idl.Greeter_TwoServer) error {
	for i := 0; i < 10; i++ {
		if err := strm.Send(&idl.Resp{Gift: "hao"}); err != nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
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
