package main

import (
	"log"
	"net"
)

func tcpListener() {
	ln, err := net.Listen("tcp", ":53")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		handleTCPConn(conn)
	}
}

func handleTCPConn(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 2048)

	for {
		n, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println(string(b[:n]))
	}
}
