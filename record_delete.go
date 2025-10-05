package provider

import (
	"context"
	"errors"
	"sync"

	"github.com/libdns/libdns"
)

func DeleteRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, deletes []libdns.Record) ([]libdns.Record, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			items, err := deleteRecords(ctx, mutex, client, zone, deletes)

			if errors.Is(err, SequentialLockerError) {
				continue
			}

			return items, err
		}
	}
}

func deleteRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, deletes []libdns.Record) ([]libdns.Record, error) {
	var unlock = lock(mutex)

	if nil != unlock {
		defer unlock()
	}

	if locker, ok := client.(SequentialLockableClient); ok {
		locker.Lock()
		defer locker.Unlock()
	}

	records, err := GetRecords(ctx, nil, client, zone)

	if err != nil {
		return nil, err
	}

	var items = make([]*libdns.RR, 0)
	var removed = make([]libdns.Record, 0)

	for _, record := range RecordIterator(&records) {
		if false == isEligibleForRemoval(&record, &deletes) {
			items = append(items, &record)
		}
	}

	if len(records) == len(items) {
		return []libdns.Record{}, nil
	}

	curr, err := client.SetDNSList(ctx, zone, items)

	if err != nil {
		return nil, err
	}

	if nil != unlock {
		unlock()
	}

	if unlock := rlock(mutex); nil != unlock {
		defer unlock()
	}

	if nil == curr {
		curr, err = GetRecords(ctx, nil, client, zone)

		if err != nil {
			return nil, err
		}
	}

	for origin, record := range RecordIterator(&records) {
		if false == IsInList(&record, &curr) && isEligibleForRemoval(&record, &deletes) {
			removed = append(removed, *origin)
		}
	}

	return removed, nil
}
