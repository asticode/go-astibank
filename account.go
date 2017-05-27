package main

import "time"

// Account represents an account
type Account struct {
	Balance    float64        `json:"balance"`
	ID         string         `json:"id"`
	Operations *OperationPool `json:"-"`
	RawBalance float64        `json:"raw_balance"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// AccountStored represents a stored account
type AccountStored struct {
	*Account
	Operations []*Operation
}

// newAccount creates a new account
func newAccount() *Account {
	return (&Account{}).init()
}

// init initializes an account
func (a *Account) init() *Account {
	a.Operations = newOperationPool()
	return a
}
