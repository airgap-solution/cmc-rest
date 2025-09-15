package main

import (
	"log"
	"net/http"
	"os"

	"github.com/airgap-solution/cmc-rest/internal/adapter/cmc"
	servicer "github.com/airgap-solution/cmc-rest/internal/adapter/http"
	"github.com/airgap-solution/cmc-rest/internal/config"
	cmcrest "github.com/airgap-solution/cmc-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/restartfu/gophig"
	"github.com/samber/lo"
)

func main() {
	configPath, _ := lo.Coalesce(os.Getenv("CONFIG_PATH"), "./config.toml")
	conf, err := loadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmcAdapter := cmc.NewAdapter(
		coinmarketcap.CurrencyBTC,
		coinmarketcap.CurrencyKAS,
		coinmarketcap.CurrencyLTC,
	)
	httpAdapter := servicer.NewAdapter(cmcAdapter)
	cmcRestServer := cmcrest.NewDefaultAPIController(httpAdapter)

	router := cmcrest.NewRouter(cmcRestServer)
	if conf.TLSEnabled {
		err = http.ListenAndServeTLS(conf.ListenAddr, conf.TLSConfig.CertificatePath, conf.TLSConfig.PrivateKeyPath, router)
	} else {
		err = http.ListenAndServe(conf.ListenAddr, router)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(configPath string) (config.Config, error) {
	defaultConfig := config.DefaultConfig()

	g := gophig.NewGophig[config.Config](configPath, gophig.TOMLMarshaler{}, 0777)
	conf, err := g.LoadConf()
	if err != nil {
		if os.IsNotExist(err) {
			err = g.SaveConf(defaultConfig)
			return defaultConfig, err
		}
		return config.Config{}, err
	}
	return conf, nil
}
