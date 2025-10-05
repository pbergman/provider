package provider

import (
	"context"
	"errors"
	"sync"

	"github.com/libdns/libdns"
)

type Domain interface {
	String() string
}

var SequentialLockerError = errors.New("update failed because remote was changed")

type Client interface {
	// GetDNSList will return all records available for given zone. This can
	// be of the opaque RR type as the provider will call Parse when available.
	GetDNSList(ctx context.Context, domain string) ([]libdns.Record, error)

	// SetDNSList will replace all the zone records with the given records.
	// If the api supports returning of set records, it will return this
	// else this can be nil, and the implementation will call GetRecords
	// to fetch a list to validate the results.
	//
	// A special case is when returning a SequentialLockerError error, then
	// the implementation will do a retry until context is done
	SetDNSList(ctx context.Context, domain string, records []*libdns.RR) ([]libdns.Record, error)
}

type SequentialLockableClient interface {
	sync.Locker
}

type ZoneAwareClient interface {
	Client
	Domains(ctx context.Context) ([]Domain, error)
}
