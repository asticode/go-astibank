package main

import (
	"encoding/json"
	"net/http"

	"bytes"
	"io/ioutil"

	"fmt"

	"encoding/csv"

	"strconv"

	"strings"

	"time"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// Vars
var (
	bytesLineSeparator = []byte("\r\n")
)

// adaptRouter adapts the router
func adaptRouter(r *httprouter.Router) {
	r.POST("/api/import", handleAPIImport)
	r.POST("/api/operations", handleAPIAddOperation)
}

// BodyOperations represents a body containing operations
type BodyOperations struct {
	Operations []BodyOperation `json:"operations"`
}

// BodyOperation represents a body containing an operation
type BodyOperation struct {
	Account   *Account   `json:"account"`
	Operation *Operation `json:"operation"`
}

// BodyPaths represents a body containing a list of paths
type BodyPaths struct {
	Paths []string `json:"paths"`
}

// processErrors processes errors
func processErrors(rw http.ResponseWriter, err *error) {
	if *err != nil {
		astilog.Error(*err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(errors.Cause(*err).Error()))
	}
}

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

		// Loop through operations
		for _, op := range ops {
			// Operation is new
			if _, err = a.Operations.One(op.ID()); err != nil {
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
	a.ID = lines[0][1]
	if a.RawBalance, err = strconv.ParseFloat(strings.Replace(lines[4][1], ",", ".", -1), 64); err != nil {
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

		// Add operation
		ops = append(ops, op)
	}
	return
}

// handleAPIAddOperation handles the /api/operations POST request
func handleAPIAddOperation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Decode input
	var bo BodyOperation
	if err = json.NewDecoder(r.Body).Decode(&bo); err != nil {
		err = errors.Wrap(err, "decoding input failed")
		return
	}

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(bo.Account.ID); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", bo.Account.ID)
		return
	}

	// Add operation
	a.Operations.Add(bo.Operation)
	rw.WriteHeader(http.StatusNoContent)
}
