package main

import (
	"time"
)

// Operation represents an operation
type Operation struct {
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	ID       int       `json:"id"`
	Label    string    `json:"label"`
	RawLabel string    `json:"raw_label"`
	Subject  string    `json:"subject"`
}
