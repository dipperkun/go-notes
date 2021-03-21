package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	// 命令行参数解析
	var server string
	var port int

	flag.StringVar(&server, "h", "127.0.0.1", "server address")
	flag.IntVar(&port, "p", 12345, "server port")
	flag.Parse()

	conn, err := net.Dial("tcp", server+":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)

		log.Println("done")
		done <- struct{}{}
	}()

	io.Copy(conn, os.Stdin)
	conn.Close()
	<-done
}
