package main

import (
	"context"
	"log"
	"net"

	helper "github.com/anmol1115/AdPhantom/internal/dns"
	Logger "github.com/anmol1115/AdPhantom/internal/logger"
	Resolver "github.com/anmol1115/AdPhantom/internal/resolver"
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

	resolver, err := Resolver.FromContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
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
		logger.Debug("Ingcoming request from ", addr.String())

		logger.Info("Parsing query")
		name, qtype, err := helper.ParseQuery(b[:n])
		if err != nil {
			logger.Error(err.Error())
			return
		}

		logger.Info("Resolving domain")
		resolvedAddr, err := resolver.Resolve(ctx, name, qtype)

		var response []byte
		if err != nil {
			logger.Error(err.Error())
			response = helper.BuildNXDomain(b[:n])
		} else {
			response = helper.BuildResponse(b[:n], resolvedAddr)
		}

		conn.WriteToUDP(response, addr)
	}
}
