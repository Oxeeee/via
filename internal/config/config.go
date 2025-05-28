package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PGSSLMode        string
	SQLitePath       string
	PGName           string
	PGUser           string
	PGPassword       string
	PGDB             string
	PublicAddr       string
	SwaggerAddr      string
	LogLevel         string
	AppName          string
	JWTRefreshSecret []byte
	JWTAccessSecret  []byte
	PGPort           int
	IsLocalRun       bool
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		AppName:          "theca",
		LogLevel:         getEnv("LOG_LEVEL", "INFO"),
		PGName:           getEnv("PG_NAME", "postgres"),
		PGUser:           getEnv("PG_USER", "postgres"),
		PGPassword:       getEnv("PG_PASSWORD", "postgres"),
		PGDB:             getEnv("PG_DB", "postgres"),
		PGPort:           getInt("PG_PORT", 5432),
		PGSSLMode:        getEnv("PG_SSL_MODE", "disable"),
		IsLocalRun:       parseBool("IS_LOCAL_RUN"),
		SQLitePath:       getEnv("SQLITE_PATH", "theca_local.db"),
		PublicAddr:       getEnv("PUBLIC_ADDR", ":8080"),
		JWTAccessSecret:  []byte(getEnv("JWT_ACCESS_SECRET", "default_access_secret")),
		JWTRefreshSecret: []byte(getEnv("JWT_REFRESH_SECRET", "default_refresh_secret")),
		SwaggerAddr:      getEnv("SWAGGER_ADDR", ":8081"),
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func parseBool(key string) bool {
	param := os.Getenv(key)
	r := false
	if param != "" {
		var err error
		r, err = strconv.ParseBool(param)
		if err != nil {
			fmt.Printf("WARN: invalid %s value (%s), defualting to false\n", key, param)
		}
	}
	return r
}

func getInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}
