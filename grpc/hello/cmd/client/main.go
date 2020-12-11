package main

import (
	"context"
	"io"
	"log"
	"sync"
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

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	cli := idl.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go one(ctx, cli, &wg)
	go two(ctx, cli, &wg)
	wg.Wait()
}

func one(ctx context.Context, cli idl.GreeterClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		resp, err := cli.One(ctx, &idl.Req{Name: "Kun"})
		if err != nil {
			log.Printf("cannot invoke one: %v", err)
		}
		log.Printf("return: %s", resp.GetGift())
	}
}

func two(ctx context.Context, cli idl.GreeterClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 20; i++ {
		resp, err := cli.One(ctx, &idl.Req{Name: "Kun"})
		if err != nil {
			log.Printf("cannot invoke one: %v", err)
		}
		log.Printf("return: %s", resp.GetGift())
	}

	strm, err := cli.Two(ctx, &idl.Req{Name: "Kun"})
	if err != nil {
		return
	}
	for {
		gift, err := strm.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v.Two = _, %v", cli, err)
		}
		log.Println(gift)
	}
}
