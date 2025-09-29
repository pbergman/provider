package transip

import (
	"context"

	"github.com/libdns/libdns"
)

type Domain interface {
	String() string
}

type Client interface {
	GetDNSList(ctx context.Context, domain string) ([]libdns.Record, error)
	SetDNSList(ctx context.Context, domain string, records []*libdns.RR) error
}

type ClientDomainProvider interface {
	Domains(ctx context.Context) ([]Domain, error)
}
