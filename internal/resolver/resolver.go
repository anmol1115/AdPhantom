package resolver

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"time"

	"github.com/anmol1115/AdPhantom/internal/config"
)

type Resolver struct {
	primary  *net.Resolver
	failover *net.Resolver
}

func New(cfg *config.Dns) *Resolver {
	return &Resolver{
		primary:  newUpstreamResolver(cfg.Upstream),
		failover: newUpstreamResolver(cfg.UpstreamFailover),
	}
}

func newUpstreamResolver(upstream string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 + time.Second}
			return d.DialContext(ctx, "udp", upstream+":53")
		},
	}
}

func (r *Resolver) Resolve(ctx context.Context, name, qtype string) ([]netip.Addr, error) {
	switch qtype {
	case "A", "AAAA":
		addrs, err := r.primary.LookupNetIP(ctx, "ip", name)
		if err != nil {
			addrs, err = r.failover.LookupNetIP(ctx, "ip", name)
		}

		return addrs, err
	default:
		return nil, errors.New("Unsupported qtype")
	}
}
