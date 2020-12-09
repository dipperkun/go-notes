package main

import (
	"context"
	"log"
	"time"

	idl "github.com/dipperkun/go-notes/grpc/hello/idl"
	"google.golang.org/grpc"
)

func main() {
	// set up a connection to the server.
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	cli := idl.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := cli.One(ctx, &idl.Req{Name: "Kun"})
	if err != nil {
		log.Fatalf("cannot invoke one: %v", err)
	}
	log.Printf("return: %s", resp.GetGift())
}
