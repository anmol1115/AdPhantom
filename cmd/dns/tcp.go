package main

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"

	Logger "github.com/anmol1115/AdPhantom/internal/logger"
	"github.com/miekg/dns"
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

	lengthBytes := make([]byte, 2)
	for {
		_, err := io.ReadFull(conn, lengthBytes)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}

			logger.Error(err.Error())
			return
		}

		msgLength := binary.BigEndian.Uint16(lengthBytes)
		msg := make([]byte, msgLength)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}

			logger.Error(err.Error())
			return
		}

		var dnsMsg dns.Msg
		err = dnsMsg.Unpack(msg)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		for _, q := range dnsMsg.Question {
			log.Println("Query: ", q.Name)
		}
	}
}
