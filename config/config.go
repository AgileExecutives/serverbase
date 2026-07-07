package config

// Config contains a minimal set of startup configuration values.
type Config struct {
	Addr        string // listen address, e.g. ":8080"
	DatabaseDsn string // optional DB connection string
}
