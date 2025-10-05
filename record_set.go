package provider

import (
	"context"
	"sync"

	"github.com/libdns/libdns"
)

func SetRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {

	var unlock = lock(mutex)

	if nil != unlock {
		defer unlock()
	}

	existing, err := GetRecords(ctx, nil, client, zone)

	if err != nil {
		return nil, err
	}

	var change = make(ChangeList, 0)

	for _, record := range RecordIterator(&existing) {

		var state = NoChange

		if nil != lookupByNameAndType(&record, &records) {
			state = Delete
		}

		change = append(change, &ChangeRecord{RR: record, State: state})
	}

	for _, item := range RecordIterator(&records) {
		change = append(change, &ChangeRecord{RR: item, State: Create})
	}

	curr, err := client.SetDNSList(ctx, zone, change)

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

	var ret = make([]libdns.Record, 0)

	for x, record := range RecordIterator(&curr) {
		if false == IsInList(&record, &existing) && nil != lookupByNameAndType(&record, &records) {
			ret = append(ret, *x)
		}
	}

	return ret, nil
}
