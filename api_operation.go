package main

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// handleAPIAddOperationToAccount handles the /api/accounts/:account_id/operations POST request
func handleAPIAddOperationToAccount(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(p.ByName("account_id")); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", p.ByName("account_id"))
		return
	}

	// Decode input
	var bo *Operation
	if err = json.NewDecoder(r.Body).Decode(&bo); err != nil {
		err = errors.Wrap(err, "decoding input failed")
		return
	}

	// Check input
	if bo.Label == "" {
		err = errors.New("Label is required")
		return
	}
	if bo.Category == "" {
		err = errors.New("Category is required")
		return
	}

	// Add operation
	a.Operations.Add(bo)
	rw.WriteHeader(http.StatusNoContent)
}

// handleAPIListOperationsOfAccount handles the /api/accounts/:account_id/operations GET request
func handleAPIListOperationsOfAccount(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(p.ByName("account_id")); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", p.ByName("account_id"))
		return
	}
	a.UpdatedAt = time.Now()

	// Write
	if err = json.NewEncoder(rw).Encode(a.Operations.All()); err != nil {
		err = errors.Wrap(err, "writing output failed")
		return
	}
}
