package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"fmt"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

type Data struct {
	Balance    float64     `json:"balance"`
	Operations []Operation `json:"operations,omitempty"`
}

type Operation struct {
	Amount      float64         `json:"amount"`
	Date        time.Time       `json:"date"`
	Description []string        `json:"description,omitempty"`
	Tags        map[string]bool `json:"tags,omitempty"`
}

type data struct {
	d    *Data
	m    *sync.Mutex
	path string
}

func newData(directory string) *data {
	return &data{
		d:    &Data{},
		m:    &sync.Mutex{},
		path: filepath.Join(directory, "data.json"),
	}
}

func (d *data) load() (err error) {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Stat file
	if _, errStat := os.Stat(d.path); errStat != nil {
		if os.IsNotExist(errStat) {
			astilog.Debugf("main: %s doesn't exist, initializing data", d.path)
		} else {
			err = errors.Wrapf(err, "main: stating %s failed", d.path)
		}
		return
	}

	// Open file
	var f *os.File
	if f, err = os.Open(d.path); err != nil {
		err = errors.Wrapf(err, "main: opening %s failed", d.path)
		return
	}
	defer f.Close()

	// Decode data
	if err = json.NewDecoder(f).Decode(d.d); err != nil {
		err = errors.Wrap(err, "main: decoding data failed")
		return
	}
	return
}

func (d *data) save() (err error) {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Create file
	var f *os.File
	if f, err = os.Create(d.path); err != nil {
		err = errors.Wrapf(err, "main: creating %s failed", d.path)
		return
	}
	defer f.Close()

	// Encode data
	if err = json.NewEncoder(f).Encode(*d.d); err != nil {
		err = errors.Wrap(err, "main: encoding data failed")
		return
	}
	return
}

func (d *data) addStatement(s pdfStatement) (err error) {
	// Lock
	d.m.Lock()
	defer d.m.Unlock()

	// Check old balance
	if len(d.d.Operations) > 0 && s.oldBalance != d.d.Balance {
		err = fmt.Errorf("old balance %f != data balance %f", s.oldBalance, d.d.Balance)
		return
	}

	// Loop through operations
	var credit, debit float64
	var os []Operation
	for _, o := range s.operations {
		credit += o.credit
		debit += o.debit
		os = append(os, Operation{
			Amount:      o.credit - o.debit,
			Date:        o.date,
			Description: o.description,
			Tags:        make(map[string]bool),
		})
	}

	// Check credit
	if credit != s.credit {
		err = fmt.Errorf("statement credit %f != computed credit %f", s.credit, credit)
		return
	}

	// Check debit
	if debit != s.debit {
		err = fmt.Errorf("statement debit %f != computed debit %f", s.debit, debit)
		return
	}

	// Add operations
	d.d.Balance = s.oldBalance + s.credit - s.debit
	d.d.Operations = append(d.d.Operations, os...)
	return
}
