package provider

import (
	"iter"

	"github.com/libdns/libdns"
)

type ChangeState uint8

const (
	NoChange ChangeState = 1 << iota
	Delete
	Create
)

type ChangeRecord struct {
	libdns.RR
	State ChangeState
}

type ChangeList []*ChangeRecord

func (c ChangeList) Iterate(state ChangeState) iter.Seq[*libdns.RR] {
	return func(yield func(*libdns.RR) bool) {
		for i, x := 0, len(c); i < x; i++ {
			if c[i].State == (c[i].State & state) {
				if false == yield(&c[i].RR) {
					return
				}
			}
		}
	}
}

func (c ChangeList) Create() []*libdns.RR {
	return c.list(Create)
}

func (c ChangeList) Deletes() []*libdns.RR {
	return c.list(Delete)
}

func (c ChangeList) GetList() []*libdns.RR {
	return c.list(Create | NoChange)
}

func (c ChangeList) list(state ChangeState) []*libdns.RR {

	var items = make([]*libdns.RR, 0)

	for record := range c.Iterate(state) {
		items = append(items, record)
	}

	return items
}
