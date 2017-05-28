package main

import (
	"fmt"
	"sync"
)

// OperationPool represents an operation pool
type OperationPool struct {
	OperationsByID map[string]*Operation
	mutex          *sync.Mutex
	OrderedIDs     []string
}

// newOperationPool creates a new operation pool
func newOperationPool() *OperationPool {
	return &OperationPool{
		OperationsByID: make(map[string]*Operation),
		mutex:          &sync.Mutex{},
	}
}

// Add adds an operation
func (p *OperationPool) Add(op *Operation) *Operation {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.OperationsByID[op.ID]; !ok {
		p.OperationsByID[op.ID] = op
		p.OrderedIDs = append(p.OrderedIDs, op.ID)
	}
	return p.OperationsByID[op.ID]
}

// All returns the operations
func (p *OperationPool) All() (os []*Operation) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	os = []*Operation{}
	for _, id := range p.OrderedIDs {
		os = append(os, p.OperationsByID[id])
	}
	return
}

// One returns the operation for a specific id
func (p *OperationPool) One(id string) (o *Operation, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var ok bool
	if o, ok = p.OperationsByID[id]; !ok {
		err = fmt.Errorf("Unknown operation id %s", id)
		return
	}
	return
}
