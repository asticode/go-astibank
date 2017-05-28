package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) {
	switch m.Name {
	case "accounts.list":
		handleMessageAccountsList(w)
	case "charts.all":
		handleMessageChartsAll(w, m)
	case "import":
		handleMessageImport(w, m)
	case "operations.add":
		handleMessageOperationsAdd(w, m)
	case "operations.list":
		handleMessageOperationsList(w, m)
	case "operations.one":
		handleMessageOperationsOne(w, m)
	case "operations.update":
		handleMessageOperationsUpdate(w, m)
	case "references.list":
		handleMessageReferencesList(w)
	}
}

// processMessageError processes the message error
func processMessageError(w *astilectron.Window, err *error) {
	if *err != nil {
		astilog.Error(*err)
		if errSend := w.Send(bootstrap.MessageOut{Name: "error", Payload: errors.Cause(*err).Error()}); errSend != nil {
			astilog.Error("Sending error message failed")
		}
	}
}
