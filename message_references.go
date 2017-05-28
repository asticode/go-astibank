package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/pkg/errors"
)

// Subjects
var (
	subject123Fleurs                = "123 Fleurs"
	subjectAccountFees              = "Account fees"
	subjectAirFrance                = "Air France"
	subjectAmazon                   = "Amazon"
	subjectAromaZone                = "Aroma Zone"
	subjectATM                      = "ATM"
	subjectButchery                 = "Butchery"
	subjectCDiscount                = "CDiscount"
	subjectCelio                    = "Celio"
	subjectDecathlon                = "Decathlon"
	subjectDeliveroo                = "Deliveroo"
	subjectEDF                      = "EDF"
	subjectEmilia                   = "Emilia"
	subjectGreenWeez                = "GreenWeez"
	subjectHerbierDeProvence        = "Herbier de Provence"
	subjectLeetchi                  = "Leetchi"
	subjectLesPrimeurs              = "Les Primeurs"
	subjectLoanInsurance            = "Loan Insurance"
	subjectMAAF                     = "MAAF"
	subjectMonoprix                 = "Monoprix"
	subjectMolotov                  = "Molotov"
	subjectOnline                   = "Online"
	subjectPharmacieDivisionLeclerc = "Pharmacie Division Leclerc"
	subjectPharmacieDuMetro         = "Pharmacie du Metro"
	subjectRATP                     = "RATP"
	subjectSelf                     = "Self"
	subjectSFR                      = "SFR"
	subjectSNCF                     = "SNCF"
	subjectTaxes                    = "Taxes"
	subjectTruffaut                 = "Truffaut"
	subjectTuaillon                 = "Tuaillon"
)

// Categories
var (
	categoryBank     = "Bank"
	categoryBread    = "Bread"
	categoryClothes  = "Clothes"
	categoryFood     = "Food"
	categoryGift     = "Gift"
	categoryHealth   = "Health"
	categoryHouse    = "House"
	categoryLoan     = "Loan"
	categoryPleasure = "Pleasure"
	categoryRent     = "Rent"
	categoryServers  = "Servers"
	categoryTaxes    = "Taxes"
	categoryUnknown  = "Unknown"
	categoryWork     = "Work"
)

// Mapping subject --> category
var mappingSubjectToCategory = map[string]string{
	subject123Fleurs:                categoryPleasure,
	subjectAccountFees:              categoryBank,
	subjectAirFrance:                categoryPleasure,
	subjectAmazon:                   categoryPleasure,
	subjectAromaZone:                categoryHealth,
	subjectATM:                      categoryFood,
	subjectButchery:                 categoryFood,
	subjectCDiscount:                categoryPleasure,
	subjectCelio:                    categoryClothes,
	subjectDecathlon:                categoryPleasure,
	subjectDeliveroo:                categoryFood,
	subjectEDF:                      categoryHouse,
	subjectEmilia:                   categoryLoan,
	subjectGreenWeez:                categoryBread,
	subjectHerbierDeProvence:        categoryFood,
	subjectLeetchi:                  categoryPleasure,
	subjectLesPrimeurs:              categoryFood,
	subjectLoanInsurance:            categoryLoan,
	subjectMAAF:                     categoryHouse,
	subjectMonoprix:                 categoryFood,
	subjectMolotov:                  categoryWork,
	subjectOnline:                   categoryServers,
	subjectPharmacieDivisionLeclerc: categoryHealth,
	subjectPharmacieDuMetro:         categoryHealth,
	subjectRATP:                     categoryWork,
	subjectSelf:                     categoryBank,
	subjectSFR:                      categoryHouse,
	subjectSNCF:                     categoryPleasure,
	subjectTaxes:                    categoryTaxes,
	subjectTruffaut:                 categoryPleasure,
	subjectTuaillon:                 categoryRent,
}

// Mapping subject --> label
var mappingSubjectToLabel = map[string]string{
	subject123Fleurs:         "Flowers",
	subjectAccountFees:       "Account fees",
	subjectATM:               "ATM Withdrawal",
	subjectButchery:          "Meat",
	subjectEDF:               "Electricity",
	subjectGreenWeez:         "Flour",
	subjectHerbierDeProvence: "Tea",
	subjectLesPrimeurs:       "Fruits & Vegetables",
	subjectLoanInsurance:     "Loan insurance",
	subjectMAAF:              "House insurance",
	subjectMolotov:           "Salary - ",
	subjectMonoprix:          "Processed food",
	subjectOnline:            "Servers",
	subjectRATP:              "Pass Navigo",
	subjectSFR:               "Internet",
	subjectTaxes:             "Taxes",
	subjectTuaillon:          "Rent - ",
}

// PayloadReferences represents the payload containing references
type PayloadReferences struct {
	Categories []string `json:"categories"`
}

// handleMessageReferencesList handles the "references.list" message
func handleMessageReferencesList(w *astilectron.Window) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "references.list", Payload: PayloadReferences{
		Categories: []string{
			categoryBank,
			categoryBread,
			categoryClothes,
			categoryFood,
			categoryGift,
			categoryHealth,
			categoryHouse,
			categoryLoan,
			categoryPleasure,
			categoryRent,
			categoryServers,
			categoryTaxes,
			categoryUnknown,
			categoryWork,
		},
	}}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}
