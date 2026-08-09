package main

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"

	Blocker "github.com/anmol1115/AdPhantom/internal/blocker"
	helper "github.com/anmol1115/AdPhantom/internal/dns"
	Logger "github.com/anmol1115/AdPhantom/internal/logger"
	Resolover "github.com/anmol1115/AdPhantom/internal/resolver"
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

	resolver, err := Resolover.FromContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fl, err := Blocker.FromContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	requestLength := make([]byte, 2)
	for {
		logger.Debug("TCP: Reading byte length of request")
		_, err = io.ReadFull(conn, requestLength)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			logger.Error(err.Error())
			return
		}

		msgLength := binary.BigEndian.Uint16(requestLength)
		if msgLength == 0 || msgLength > 4096 {
			logger.Error("TCP: Invalid length of request message")
			return
		}

		logger.Debug("TCP: Reading request message")
		msg := make([]byte, msgLength)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			logger.Error(err.Error())
			return
		}

		logger.Debug("TCP: Parsing query")
		name, qtype, err := helper.ParseQuery(msg)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		var tmp []byte
		rule := fl.Match(name)
		if rule.Type == Blocker.ExactAllow || rule.Type == Blocker.WildcardAllow {
			logger.Debug("TCP: Resolving domain")
			resolvedAddr, err := resolver.Resolve(ctx, name, qtype)

			if err != nil {
				logger.Error(err.Error())
				tmp = helper.BuildNXDomain(msg)
			} else {
				tmp = helper.BuildResponse(msg, resolvedAddr)
			}
		} else {
			tmp = helper.BuildNXDomain(msg)
		}

		logger.Debug("TCP: Sending response")
		var response []byte
		response = binary.BigEndian.AppendUint16(response, uint16(len(tmp)))
		response = append(response, tmp...)

		conn.Write(response)
	}
}
