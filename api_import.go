package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// handleAPIImport handles the /api/import POST request
func handleAPIImport(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Decode input
	var bp BodyPaths
	if err = json.NewDecoder(r.Body).Decode(&bp); err != nil {
		err = errors.Wrap(err, "decoding input failed")
		return
	}

	// Loop in paths
	var bo = BodyOperations{}
	for _, p := range bp.Paths {
		// Parse bank statement
		var a *Account
		var ops []*Operation
		if a, ops, err = parseBankStatement(p); err != nil {
			err = errors.Wrapf(err, "parsing bank statement %s failed", p)
			return
		}

		// Set account
		a = data.Accounts.Set(a)

		// Add new operations
		for _, op := range ops {
			if _, err = a.Operations.One(op.ID); err != nil {
				bo.Operations = append(bo.Operations, BodyOperation{Account: a, Operation: op})
			}
		}
	}

	// Write
	if err = json.NewEncoder(rw).Encode(bo); err != nil {
		err = errors.Wrap(err, "writing output failed")
		return
	}
}

// parseBankStatement parses a bank statement
func parseBankStatement(path string) (a *Account, ops []*Operation, err error) {
	// Log
	astilog.Debugf("Parsing bank statement %s", path)

	// Check file extension
	if filepath.Ext(path) != ".csv" {
		err = fmt.Errorf("invalid extension for %s", path)
		return
	}

	// Open file
	var b []byte
	if b, err = ioutil.ReadFile(path); err != nil {
		err = errors.Wrapf(err, "opening %s failed", path)
		return
	}

	// Split header from body
	var items = bytes.Split(b, append(bytesLineSeparator, bytesLineSeparator...))
	if len(items) == 1 {
		err = fmt.Errorf("no body detected in content %s", b)
		return
	}

	// Build header csv reader
	var hr = csv.NewReader(bytes.NewReader(items[0]))
	hr.Comma = ';'
	hr.FieldsPerRecord = 2

	// Read header lines
	var lines [][]string
	if lines, err = hr.ReadAll(); err != nil {
		err = errors.Wrapf(err, "reading header lines of %s failed", path)
		return
	}
	if len(lines) < 6 {
		err = fmt.Errorf("not enough lines in header %s", items[0])
		return
	}

	// Parse account fields
	a = newAccount()
	a.ID = fmt.Sprintf("%s %s", lines[1][1], lines[0][1])
	if a.Balance, err = strconv.ParseFloat(strings.Replace(lines[4][1], ",", ".", -1), 64); err != nil {
		err = fmt.Errorf("%s is not a valid float", lines[4][1])
		return
	}

	// Build body csv reader
	var br = csv.NewReader(bytes.NewReader(items[1]))
	br.Comma = ';'
	br.FieldsPerRecord = 4

	// Read body lines
	br.Read()
	if lines, err = br.ReadAll(); err != nil {
		err = errors.Wrapf(err, "reading body lines of %s failed", path)
		return
	}

	// Loop through lines
	for i := len(lines) - 1; i >= 0; i-- {
		// Init
		var op = &Operation{RawLabel: lines[i][1]}

		// Parse date
		if op.Date, err = time.Parse("02/01/2006", lines[i][0]); err != nil {
			err = fmt.Errorf("%s is not a valid date", lines[i][0])
			return
		}

		// Parse amount
		if op.Amount, err = strconv.ParseFloat(strings.Replace(lines[i][2], ",", ".", -1), 64); err != nil {
			err = fmt.Errorf("%s is not a valid float", lines[i][2])
			return
		}

		// Update account balance
		a.Balance -= op.Amount

		// Set ID
		op.ID = fmt.Sprintf("%s.%s.%f", op.Date, op.RawLabel, op.Amount)

		// Add operation
		ops = append(ops, op)
	}
	return
}
