package main

import (
	"fmt"
	"sync"
)

// OperationPool represents an operation pool
type OperationPool struct {
	Counter        int
	OperationsByID map[int]*Operation
	mutex          *sync.Mutex
	OrderedIDs     []int
}

// newOperationPool creates a new operation pool
func newOperationPool() *OperationPool {
	return &OperationPool{
		OperationsByID: make(map[int]*Operation),
		mutex:          &sync.Mutex{},
	}
}

// Add adds an operation
func (p *OperationPool) Add(op *Operation) *Operation {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Counter++
	op.ID = p.Counter
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
func (p *OperationPool) One(id int) (o *Operation, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var ok bool
	if o, ok = p.OperationsByID[id]; !ok {
		err = fmt.Errorf("Unknown operation id %d", id)
		return
	}
	return
}
