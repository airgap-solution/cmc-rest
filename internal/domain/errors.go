package domain

import "errors"

var (
	ErrRateNotFound = errors.New("rate not found")
	ErrInvalidFiat  = errors.New("invalid fiat currency")
)
