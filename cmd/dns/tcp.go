package main

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"

	Logger "github.com/anmol1115/AdPhantom/internal/logger"
	"github.com/miekg/dns"
)

func tcpListener(ctx context.Context) {
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("TCP: Starting listener on port 53")
	ln, err := net.Listen("tcp", ":53")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		logger.Debug("TCP: Closing listener")
		ln.Close()
	}()

	for {
		logger.Debug("TCP: Waiting for connection")
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				logger.Debug("TCP: Exiting listener loop")
				return
			}
			logger.Error(err.Error())
			continue
		}
		go handleTCPConn(ctx, conn)
	}
}

func handleTCPConn(ctx context.Context, conn net.Conn) {
	logger, err := Logger.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("TCP: Connection established from: " + conn.RemoteAddr().String())
	defer func() {
		logger.Debug("TCP: Connection closed from: " + conn.RemoteAddr().String())
		conn.Close()
	}()

	lengthBytes := make([]byte, 2)
	var dnsMsg dns.Msg

	for {
		logger.Debug("TCP: Reading byte length of request")
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err := io.ReadFull(conn, lengthBytes)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}

			logger.Error(err.Error())
			return
		}

		logger.Debug("TCP: Reading request message")
		msgLength := binary.BigEndian.Uint16(lengthBytes)
		if msgLength == 0 || msgLength > 4096 {
			logger.Error("TCP: Invalid length of message received")
			return
		}

		msg := make([]byte, msgLength)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}

			logger.Error(err.Error())
			return
		}

		logger.Debug("TCP: Parsing request")
		dnsMsg = dns.Msg{}
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
