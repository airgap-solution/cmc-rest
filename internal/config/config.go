package config

type Config struct {
	ListenAddr string `toml:"listen_addr"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddr: ":8080",
	}
}
