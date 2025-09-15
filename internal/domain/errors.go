package domain

import "errors"

var (
	ErrRateNotFound     = errors.New("rate not found")
	ErrInvalidFiat      = errors.New("invalid fiat currency")
	ErrInvalidTimeFrame = errors.New("invalid time frame")
)
