package main

import (
	"fmt"
	"time"
)

// Operation represents an operation
type Operation struct {
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Label    string    `json:"label"`
	RawLabel string    `json:"raw_label"`
}

// ID builds the operation's ID
func (o Operation) ID() string {
	return fmt.Sprintf("%s.%s.%f", o.Date, o.RawLabel, o.Amount)
}
