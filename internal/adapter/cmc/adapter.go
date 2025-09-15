package cmc

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/airgap-solution/cmc-rest/internal/domain"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/restartfu/coinmarketcap/coinmarketcap/fiat"
)

type Adapter struct {
	subscriber    *coinmarketcap.Subscriber
	rateConverter *fiat.DefaultRateConverter
	cache         map[string]rateUpdate
	mu            sync.RWMutex
}

func NewAdapter(currencies ...coinmarketcap.Currency) *Adapter {
	a := &Adapter{
		subscriber:    coinmarketcap.NewSubscriber(fiat.USD, nil),
		rateConverter: fiat.NewDefaultRateConverter(nil),
		cache:         make(map[string]rateUpdate),
	}

	a.subscriber.Subscribe(currencies...)
	go a.rateConverter.Start(time.Minute)
	go a.pollLoop(currencies)

	return a
}

func (a *Adapter) pollLoop(currencies []coinmarketcap.Currency) {
	for _, curr := range currencies {
		go func(c coinmarketcap.Currency) {
			for {
				details, err := a.subscriber.Poll(c)
				if err != nil {
					log.Printf("poll error for %s: %v", c, err)
					a.subscriber.Revive()
					time.Sleep(5 * time.Second) // avoid tight retry loop
					continue
				}

				change24 := details.Price - (details.Price / (1 + details.Price24H/100))
				a.mu.Lock()
				a.cache[strings.ToUpper(c.String())] = rateUpdate{
					rate:     details.Price,
					updateAt: time.Now(),
					change24: change24,
				}
				a.mu.Unlock()
			}
		}(curr)
	}
}

func (a *Adapter) GetCryptoFiatRate(crypto, fiatSymbol string) (float64, time.Time, error) {
	crypto = strings.ToUpper(crypto)
	fiatSymbol = strings.ToUpper(fiatSymbol)

	fiatCurrency, ok := fiat.BySymbol(fiatSymbol)
	if !ok {
		return 0, time.Time{}, fmt.Errorf("%s: %w", fiatSymbol, domain.ErrInvalidFiat)
	}

	a.mu.RLock()
	priceUSD, ok := a.cache[crypto]
	a.mu.RUnlock()
	if !ok {
		return 0, time.Time{}, fmt.Errorf("%s: %w", crypto, domain.ErrRateNotFound)
	}

	rate, err := a.rateConverter.ConvertRate(fiat.USD, fiatCurrency)
	if err != nil {
		return 0, time.Time{}, err
	}

	return priceUSD.rate * rate, priceUSD.updateAt, nil
}
func (a *Adapter) GetCryptoFiatChange(crypto, fiatSymbol, timeFrame string) (float64, error) {
	if timeFrame != "24h" {
		return 0, fmt.Errorf("%s: %w", timeFrame, domain.ErrInvalidTimeFrame)
	}

	crypto = strings.ToUpper(crypto)
	fiatSymbol = strings.ToUpper(fiatSymbol)

	fiatCurrency, ok := fiat.BySymbol(fiatSymbol)
	if !ok {
		return 0, fmt.Errorf("%s: %w", fiatSymbol, domain.ErrInvalidFiat)
	}

	a.mu.RLock()
	entry, ok := a.cache[crypto]
	a.mu.RUnlock()
	if !ok {
		return 0, fmt.Errorf("%s: %w", crypto, domain.ErrRateNotFound)
	}

	rate, err := a.rateConverter.ConvertRate(fiat.USD, fiatCurrency)
	if err != nil {
		return 0, err
	}

	return entry.change24 * rate, nil
}

type rateUpdate struct {
	rate     float64
	updateAt time.Time
	change24 float64
}
