package main

import (
	"net/http"

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
	r.GET("/api/accounts", handleAPIListAccounts)
	r.POST("/api/accounts/:account_id/operations", handleAPIAddOperationToAccount)
	r.GET("/api/accounts/:account_id/operations", handleAPIListOperationsOfAccount)
	r.POST("/api/import", handleAPIImport)
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
