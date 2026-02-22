package main

import (
	"context"
	"io"
	"log"
	"net"

	Logger "github.com/anmol1115/AdPhantom/internal/logger"
)

func tcpListener(ctx context.Context) {
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("Starting TCP Listener on port 53")
	ln, err := net.Listen("tcp", ":53")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		logger.Debug("Closing TCP Listener")
		ln.Close()
	}()

	for {
		logger.Debug("Waiting for TCP connection")
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				logger.Debug("Closing TCP Listener")
				return
			}
			logger.Error(err.Error())
			continue
		}
		logger.Debug("TCP connection established")
		go handleTCPConn(ctx, conn)
	}
}

func handleTCPConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2048)
	for {
		n, err := conn.Read(b)
		if err == io.EOF {
			return
		}
		if err != nil {
			logger.Error(err.Error())
			return
		}
		log.Println(string(b[:n]))
	}
}
