package config

type Config struct {
	ListenAddr string `toml:"listen_addr"`
	TLSEnabled bool   `toml:"tls_enabled"`
	TLSConfig  struct {
		CertificatePath string `toml:"certificate_path"`
		PrivateKeyPath  string `toml:"private_key_path"`
	} `toml:"tls_config"`
}

func DefaultConfig() Config {
	cfg := Config{
		ListenAddr: "restartfu.com:8765",
		TLSEnabled: true,
	}

	cfg.TLSConfig.CertificatePath = "/etc/letsencrypt/live/restartfu.com/fullchain.pem"
	cfg.TLSConfig.PrivateKeyPath = "/etc/letsencrypt/live/restartfu.com/privkey.pem"
	return cfg
}
