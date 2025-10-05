package provider

import (
	"context"
	"errors"
	"sync"

	"github.com/libdns/libdns"
)

func SetRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:

			items, err := setRecords(ctx, mutex, client, zone, records)

			if errors.Is(err, SequentialLockerError) {
				continue
			}

			return items, err
		}
	}
}

func setRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var unlock = lock(mutex)

	if nil != unlock {
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

	var set = make([]*libdns.RR, 0)
	var ret = make([]libdns.Record, 0)

	for _, record := range RecordIterator(&existing) {
		if nil == lookupByNameAndType(&record, &records) {
			set = append(set, &record)
		}
	}

	for _, item := range RecordIterator(&records) {
		set = append(set, &item)
	}

	curr, err := client.SetDNSList(ctx, zone, set)

	if err != nil {
		return nil, err
	}

	if nil != unlock {
		unlock()
	}

	if unlock := rlock(mutex); nil != unlock {
		defer unlock()
	}

	if curr == nil {
		curr, err = GetRecords(ctx, nil, client, zone)

		if err != nil {
			return nil, err
		}
	}

	for x, record := range RecordIterator(&curr) {
		if false == IsInList(&record, &existing) && nil != lookupByNameAndType(&record, &records) {
			ret = append(ret, *x)
		}
	}

	return ret, nil
}
