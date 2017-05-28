package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/pkg/errors"
)

// handleMessageAccountsList handles the "accounts.list" message
func handleMessageAccountsList(w *astilectron.Window) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "accounts.list", Payload: data.Accounts.All()}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
