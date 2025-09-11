package handler

import (
	"context"

	"github.com/airgap-solution/cmc-rest/internal/ports"
	cmcrest "github.com/airgap-solution/cmc-rest/openapi/servergen/go"
)

type Adapter struct {
	cmcAdapter ports.CMC
}

func NewAdapter(cmcAdapter ports.CMC) *Adapter {
	return &Adapter{cmcAdapter: cmcAdapter}
}

func (a *Adapter) V1RateCurrencyFiatGet(_ context.Context, from string, to string) (cmcrest.ImplResponse, error) {
	rate, updatedAt, err := a.cmcAdapter.GetCryptoFiatRate(from, to)
	if err != nil {
		return cmcrest.Response(500, nil), err
	}
	return cmcrest.Response(200, cmcrest.GetRateResponse{
		Rate:      rate,
		UpdatedAt: updatedAt,
	}), nil
}
