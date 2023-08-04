package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"rrpc/message"
	"rrpc/server"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	gobServer := message.NewGobMessage(conn)
	for i := 0; i < 5; i++ {
		h := &message.Header{
			Version:  1,
			SeqNum:   uint64(i),
			Error:    "",
			AuthType: 0,
		}
		_ = gobServer.WriteMessage(h, fmt.Sprintf("server rpc %d", h.SeqNum))
		_ = gobServer.ReadHeader(h)
		var reply string
		_ = gobServer.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
