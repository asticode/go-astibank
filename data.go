package main

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Data represents data
type Data struct {
	Accounts *accountPool
	path     string
}

// dataPath returns the data path
func dataPath(baseDirPath string) string {
	return filepath.Join(baseDirPath, "data.bin")
}

// NewData creates new data
func NewData(baseDirPath string) (d *Data, err error) {
	// Init
	d = &Data{
		Accounts: newAccountPool(),
		path:     dataPath(baseDirPath),
	}

	// Open data file
	var f *os.File
	if f, err = os.Open(d.path); os.IsNotExist(err) {
		astilog.Debugf("%s doesn't exist, working with new data", d.path)
		err = nil
		return
	} else if err != nil {
		err = errors.Wrapf(err, "stating %s failed", d.path)
		return
	}

	// Decode data
	astilog.Debugf("Importing data from %s", d.path)
	var ass []AccountStored
	if err = gob.NewDecoder(f).Decode(&ass); err != nil {
		err = errors.Wrapf(err, "decoding %s failed", d.path)
		return
	}

	// Loop through accounts
	for _, as := range ass {
		// Set account
		var a = d.Accounts.Set(as.init())

		// Loop through operations
		for _, o := range as.Operations {
			a.Operations.Add(o)
		}
	}
	return
}

// Close closes the data properly
func (d *Data) Close() (err error) {
	// Create file
	var f *os.File
	if f, err = os.Create(d.path); err != nil {
		err = errors.Wrapf(err, "creating %s failed", d.path)
		return
	}

	// Build data
	var ass []AccountStored
	for _, a := range d.Accounts.All() {
		var as = AccountStored{Account: a}
		for _, o := range a.Operations.All() {
			as.Operations = append(as.Operations, o)
		}
		ass = append(ass, as)
	}

	// Encode data
	astilog.Debugf("Exporting data to %s", d.path)
	if err = gob.NewEncoder(f).Encode(ass); err != nil {
		err = errors.Wrapf(err, "encoding %s failed", d.path)
		return
	}
	return
}
