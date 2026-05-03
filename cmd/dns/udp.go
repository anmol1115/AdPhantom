package main

import (
	"context"
	"log"
	"net"

	Logger "github.com/anmol1115/AdPhantom/internal/logger"
)

func udpListener(ctx context.Context) {
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("Creating UDP address")
	udpAddr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("Starting UDP listener on port 53")
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer udpConn.Close()

	go handleUDPConn(ctx, udpConn)

	<-ctx.Done()
	logger.Debug("Closing UDP listener")
}

func handleUDPConn(ctx context.Context, conn *net.UDPConn) {
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			if ctx.Err() != nil {
				logger.Debug("Closing UDP Listener")
				return
			}
			logger.Error(err.Error())
			continue
		}
		logger.Debug(addr.String(), string(b[:n]))
	}
}
