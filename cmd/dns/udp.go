package main

import (
	"log"
	"net"
)

func UdpListener() {
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
