package main

import (
	"encoding/json"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/pkg/errors"
)

// handleMessageOperationsAdd handles the "operations.add" message
func handleMessageOperationsAdd(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal
	var po PayloadOperation
	if err = json.Unmarshal(m.Payload, &po); err != nil {
		err = errors.Wrapf(err, "unmarshaling %s failed", m.Payload)
		return
	}

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(po.Account.ID); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", po.Account.ID)
		return
	}

	// Check input
	if po.Operation.Subject == "" {
		err = errors.New("Subject is required")
		return
	}
	if po.Operation.Category == "" {
		err = errors.New("Category is required")
		return
	}
	if po.Operation.Label == "" {
		err = errors.New("Label is required")
		return
	}

	// Add operation
	a.Operations.Add(po.Operation)

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "operations.add"}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}

// handleMessageOperationsList handles the "operations.list" message
func handleMessageOperationsList(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal
	var pa string
	if err = json.Unmarshal(m.Payload, &pa); err != nil {
		err = errors.Wrapf(err, "unmarshaling %s failed", m.Payload)
		return
	}

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(pa); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", pa)
		return
	}
	a.UpdatedAt = time.Now()

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "operations.list", Payload: a.Operations.All()}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
