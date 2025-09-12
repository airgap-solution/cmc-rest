package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/airgap-solution/cmc-rest/internal/adapter/cmc"
	servicer "github.com/airgap-solution/cmc-rest/internal/adapter/http"
	cmcrest "github.com/airgap-solution/cmc-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/samber/lo"
)

func localIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				return ip4.String(), nil
			}
		}
	}
	return "", nil
}

func main() {
	cmcAdapter := cmc.NewAdapter(
		coinmarketcap.CurrencyBTC,
		coinmarketcap.CurrencyKAS,
		coinmarketcap.CurrencyLTC,
	)
	httpAdapter := servicer.NewAdapter(cmcAdapter)
	cmcRestServer := cmcrest.NewDefaultAPIController(httpAdapter)

	ip := lo.Must(localIP())
	log.Printf("Server listening at %s:8083", ip)

	router := cmcrest.NewRouter(cmcRestServer)
	err := http.ListenAndServe(fmt.Sprintf("%s:8083", ip), router)
	if err != nil {
		log.Fatalln(err)
	}
}
