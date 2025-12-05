package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  	string
	GRPCPort 	string

	DBHost  	string
	DBPort  	string
	DBUser  	string
	DBPass  	string
	DBName  	string
	DBSSL   	string
	DBShardsRaw string

	GinMode 	string
	AppEnv  	string
}

type Loader interface {
	Load() (*Config, error)
}

type EnvLoader struct{}

func NewEnvLoader() Loader {
	return &EnvLoader{}
}

func (l *EnvLoader) Load() (*Config, error) {
	_ = godotenv.Load(".env.dev")

	cfg := &Config{
		AppEnv:  getEnvOrDefault("APP_ENV", "dev"),
		AppPort: getEnvOrDefault("APP_PORT", "8080"),
		GRPCPort: getEnvOrDefault("GRPC_PORT", "50051"),

		DBHost: getEnvOrDefault("DB_HOST", "localhost"),
		DBPort: getEnvOrDefault("DB_PORT", "5432"),
		DBUser: getEnvOrDefault("DB_USER", "postgres"),
		DBPass: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName: getEnvOrDefault("DB_NAME", "notes"),
		DBSSL:  getEnvOrDefault("DB_SSLMODE", "disable"),

		DBShardsRaw: getEnvOrDefault("DB_SHARDS", "db,db-2"),

		GinMode: getEnvOrDefault("GIN_MODE", "release"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func MustLoad(loader Loader) *Config {
	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func (c *Config) DBShardHosts() []string {
    if c.DBShardsRaw == "" {
        return nil
    }
    parts := strings.Split(c.DBShardsRaw, ",")
    var res []string
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" {
            res = append(res, p)
        }
    }
    return res
}
