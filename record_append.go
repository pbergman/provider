package provider

import (
	"context"
	"errors"
	"sync"

	"github.com/libdns/libdns"
)

// AppendRecords appends new records to the existing ones without performing validation.
// The assumption is that when the full list is returned to the provider,
// the provider will handle any necessary validation and fail if problems are found.
func AppendRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			items, err := appendRecords(ctx, mutex, client, zone, records)

			if errors.Is(err, SequentialLockerError) {
				continue
			}

			return items, err
		}
	}
}

func appendRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {

	if unlock := lock(mutex); unlock != nil {
		defer unlock()
	}

	if locker, ok := client.(SequentialLockableClient); ok {
		locker.Lock()
		defer locker.Unlock()
	}

	existing, err := GetRecords(ctx, nil, client, zone)

	if err != nil {
		return nil, err
	}

	var items = make([]*libdns.RR, 0, len(existing)+len(records))
	var newList = append(existing, records...)

	for _, record := range RecordIterator(&newList) {
		items = append(items, &record)
	}

	curr, err := client.SetDNSList(ctx, zone, items)

	if err != nil {
		return nil, err
	}

	if nil == curr {
		curr, err = GetRecords(ctx, nil, client, zone)

		if err != nil {
			return nil, err
		}
	}

	var ret = make([]libdns.Record, 0)

	for origin, record := range RecordIterator(&curr) {
		if false == IsInList(&record, &existing) {
			ret = append(ret, *origin)
		}
	}

	return ret, nil
}
