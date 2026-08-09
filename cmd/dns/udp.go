package main

import (
	"context"
	"log"
	"net"

	Blocker "github.com/anmol1115/AdPhantom/internal/blocker"
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

	fl, err := Blocker.FromContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	b := make([]byte, 2048)
	for {
		logger.Debug("UDP: Reading byte length of request")
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			if ctx.Err() != nil {
				logger.Debug("UDP: Closing UDP Listener")
				return
			}
			logger.Error(err.Error())
			continue
		}
		logger.Debug("UDP: Ingcoming request from ", addr.String())

		logger.Debug("UDP: Parsing query")
		name, qtype, err := helper.ParseQuery(b[:n])
		if err != nil {
			logger.Error(err.Error())
			return
		}

		var response []byte
		rule := fl.Match(name)
		if rule.Type == Blocker.ExactAllow || rule.Type == Blocker.WildcardAllow {
			logger.Debug("UDP: Resolving domain")
			resolvedAddr, err := resolver.Resolve(ctx, name, qtype)

			if err != nil {
				logger.Error(err.Error())
				response = helper.BuildNXDomain(b[:n])
			} else {
				response = helper.BuildResponse(b[:n], resolvedAddr)
			}
		} else {
			response = helper.BuildNXDomain(b[:n])
		}

		logger.Debug("UDP: Sending response")
		conn.WriteToUDP(response, addr)
	}
}
