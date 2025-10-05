package provider

import (
	"context"
	"sync"

	"github.com/libdns/libdns"
)

// AppendRecords appends new records to the existing ones without performing validation.
// The assumption is that when the full list is returned to the provider,
// the provider will handle any necessary validation and fail if problems are found.
func AppendRecords(ctx context.Context, mutex sync.Locker, client Client, zone string, records []libdns.Record) ([]libdns.Record, error) {

	if unlock := lock(mutex); unlock != nil {
		defer unlock()
	}

	existing, err := GetRecords(ctx, nil, client, zone)

	if err != nil {
		return nil, err
	}

	var change = make(ChangeList, 0, len(existing)+len(records))

	for i, c := 0, len(existing); i < c; i++ {
		change = append(change, &ChangeRecord{RR: existing[i].RR(), State: NoChange})
	}

	for i, c := 0, len(records); i < c; i++ {
		change = append(change, &ChangeRecord{RR: records[i].RR(), State: Create})
	}

	items, err := client.SetDNSList(ctx, zone, change)

	if err != nil {
		return nil, err
	}

	if nil == items {
		items, err = GetRecords(ctx, nil, client, zone)

		if err != nil {
			return nil, err
		}
	}

	var ret = make([]libdns.Record, 0)

	for origin, record := range RecordIterator(&items) {
		if false == IsInList(&record, &existing) {
			ret = append(ret, *origin)
		}
	}

	return ret, nil
}
