package ports

import "time"

type CMC interface {
	GetCryptoFiatRate(crypto, fiat string) (float64, time.Time, error)
	GetCryptoFiatChange(crypto, fiat, timeFrame string) (float64, error)
}
