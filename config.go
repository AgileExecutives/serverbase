package serverbase

import "os"

// Config contains minimal server configuration.
type Config struct {
	Port           string
	MigrateOnStart bool
}

// LoadConfigFromEnv loads minimal config from environment variables.
// - SERVER_PORT: port to bind (default :8085)
// - MIGRATE_ON_START: if set to "true" will run Module.Migrate() before start
func LoadConfigFromEnv() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8085"
	}
	migrate := false
	if os.Getenv("MIGRATE_ON_START") == "true" {
		migrate = true
	}
	return &Config{Port: port, MigrateOnStart: migrate}
}
