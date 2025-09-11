package main

import (
	"log"
	"net/http"

	"github.com/airgap-solution/cmc-rest/internal/adapter/cmc"
	httpHandler "github.com/airgap-solution/cmc-rest/internal/adapter/http"
	cmcrest "github.com/airgap-solution/cmc-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
)

func main() {
	cmcAdapter := cmc.NewAdapter(
		coinmarketcap.CurrencyBTC,
		coinmarketcap.CurrencyKAS,
		coinmarketcap.CurrencyLTC,
	)
	httpAdapter := httpHandler.NewAdapter(cmcAdapter)
	cmcRestServer := cmcrest.NewDefaultAPIController(httpAdapter)

	router := cmcrest.NewRouter(cmcRestServer)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err)
	}
}
