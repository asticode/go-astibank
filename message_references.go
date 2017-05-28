package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/pkg/errors"
)

// Subjects
var (
	subjectATMBanquePostale = "ATM Banque Postale"
	subjectDecathlon        = "Decathlon"
	subjectEDF              = "EDF"
	subjectLesPrimeurs      = "Les Primeurs"
)

// Categories
var (
	categoryElectricity = "Electricity"
	categoryFood        = "Food"
	categoryPleasure    = "Pleasure"
)

// Mapping subject --> category
var mappingSubjectToCategory = map[string]string{
	subjectATMBanquePostale: categoryFood,
	subjectDecathlon:        categoryPleasure,
	subjectEDF:              categoryElectricity,
	subjectLesPrimeurs:      categoryFood,
}

// Mapping subject --> label
var mappingSubjectToLabel = map[string]string{
	subjectATMBanquePostale: "ATM withdrawal",
	subjectEDF:              "Electricity",
	subjectLesPrimeurs:      "Fruits & Vegetables",
}

// PayloadReferences represents the payload containing references
type PayloadReferences struct {
	Categories []string `json:"categories"`
	Subjects   []string `json:"subjects"`
}

// handleMessageReferencesList handles the "references.list" message
func handleMessageReferencesList(w *astilectron.Window) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "references.list", Payload: PayloadReferences{
		Categories: []string{
			categoryElectricity,
			categoryFood,
			categoryPleasure,
		},
		Subjects: []string{
			subjectATMBanquePostale,
			subjectDecathlon,
			subjectEDF,
			subjectLesPrimeurs,
		},
	}}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
