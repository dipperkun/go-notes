package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type msg struct {
	owner   uint64
	content string
}

type user struct {
	id      uint64
	addr    string
	enterAt time.Time
	msgChan chan msg
}

func (u *user) String() string {
	return strconv.Itoa(int(u.id)) + "[" + u.addr + "]"
}

var userID uint64

func id() uint64 {
	return atomic.AddUint64(&userID, 1)
}

var (
	// new user
	enteringChan = make(chan *user)
	// user leave
	leavingChan = make(chan *user)
	// msg
	msgChan = make(chan msg, 20)
)

func main() {
	var port int
	flag.IntVar(&port, "p", 12345, "listen port")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	// 广播消息
	go broadcast()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go chat(conn)
	}
}

func broadcast() {
	users := make(map[*user]struct{})

	for {
		select {
		case u := <-enteringChan:
			users[u] = struct{}{}
		case u := <-leavingChan:
			delete(users, u)
			close(u.msgChan) // warning
		case msg := <-msgChan:
			for u := range users {
				if u.id == msg.owner {
					continue
				}
				u.msgChan <- msg
			}
		}
	}
}

func chat(conn net.Conn) {
	defer conn.Close()

	// new user comes in
	u := &user{
		id:      id(),
		addr:    conn.RemoteAddr().String(),
		enterAt: time.Now(),
		msgChan: make(chan msg, 100),
	}

	// send msg
	go sendMsg(conn, u.msgChan)

	u.msgChan <- msg{owner: 0, content: "welcome, " + u.String()}
	msgChan <- msg{owner: 0, content: "user: " + u.String() + " enter"}

	enteringChan <- u

	var active = make(chan struct{})
	go func() {
		duration := time.Minute
		timer := time.NewTimer(duration)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-active:
				timer.Reset(duration)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msgChan <- msg{owner: u.id, content: input.Text()}

		active <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("read:", err)
	}

	leavingChan <- u
	msgChan <- msg{owner: 0, content: u.String() + "leave"}
}

func sendMsg(conn net.Conn, ch <-chan msg) {
	for m := range ch {
		fmt.Fprintln(conn, strconv.Itoa(int(m.owner))+": "+m.content)
	}
}
