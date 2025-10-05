package provider

import (
	"context"

	"github.com/libdns/libdns"
)

type Domain interface {
	String() string
}

type Client interface {
	// GetDNSList will return all records available for given zone. This can
	// be of the opaque RR type as the provider will call Parse when available
	GetDNSList(ctx context.Context, domain string) ([]libdns.Record, error)

	// SetDNSList should walk through the change list and update the records
	// depending on there state. This can be one of NoChange, Delete or Create
	// The client coder can be something like:
	//
	// for remove := change.Iterate(Delete) {
	//    ... //remove record
	// }
	//
	// for create := change.Iterate(Create) {
	//    ... //create record
	// }
	//
	// Or call GetList to generate a new list which can be used to replace the
	// whole record list at once
	//
	// client.setList(ctx, change.GetList())
	SetDNSList(ctx context.Context, domain string, change ChangeList) ([]libdns.Record, error)
}

type ZoneAwareClient interface {
	Client
	Domains(ctx context.Context) ([]Domain, error)
}
