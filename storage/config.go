package storage

type Config struct {
	Provider string `koanf:"provider"`
	DSN      string `koanf:"dsn"`
}
