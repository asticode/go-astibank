package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

// BodyReferences represents the body containing references
type BodyReferences struct {
	Categories []string `json:"categories"`
	Subjects   []string `json:"subjects"`
}

// handleAPIListReferences handles the /api/references GET request
func handleAPIListReferences(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process errors
	var err error
	defer processErrors(rw, &err)

	// Write
	if err = json.NewEncoder(rw).Encode(BodyReferences{
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
	}); err != nil {
		err = errors.Wrap(err, "writing output failed")
		return
	}
}
