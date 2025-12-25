package main

import (
	"log"
	"net"
)

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

func handleUDPConn(conn *net.UDPConn) {
	b := make([]byte, 2048)

	for {
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println(addr, string(b[:n]))
	}
}

func main() {
	// TCP part
	go func() {
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
	}()

	// UDP part
	go func() {
		udpAddr, err := net.ResolveUDPAddr("udp", ":53")
		if err != nil {
			log.Fatal(err)
		}
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer udpConn.Close()

		handleUDPConn(udpConn)
	}()
}
