package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// handleAPIListAccounts handles the /api/accounts GET request
func handleAPIListAccounts(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Write
	if err = json.NewEncoder(rw).Encode(data.Accounts.All()); err != nil {
		err = errors.Wrap(err, "writing output failed")
		return
	}
}
