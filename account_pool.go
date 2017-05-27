package main

import (
	"fmt"
	"sort"
	"sync"
)

// accountPool represents an account pool
type accountPool struct {
	accountsByID map[string]*Account
	mutex        *sync.Mutex
	orderedIDs   []string
}

// newAccountPool creates a new account pool
func newAccountPool() *accountPool {
	return &accountPool{
		accountsByID: make(map[string]*Account),
		mutex:        &sync.Mutex{},
	}
}

// All returns the accounts
func (p *accountPool) All() (as []*Account) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, id := range p.orderedIDs {
		as = append(as, p.accountsByID[id])
	}
	return
}

// One returns the account for a specific id
func (p *accountPool) One(id string) (a *Account, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var ok bool
	if a, ok = p.accountsByID[id]; !ok {
		err = fmt.Errorf("Unknown account id %s", id)
		return
	}
	return
}

// Set sets an account
func (p *accountPool) Set(a *Account) *Account {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.accountsByID[a.ID]; !ok {
		p.accountsByID[a.ID] = a
		p.orderedIDs = append(p.orderedIDs, a.ID)
		sort.Strings(p.orderedIDs)
	} else {
		p.accountsByID[a.ID].RawBalance = a.RawBalance
	}
	return p.accountsByID[a.ID]
}
