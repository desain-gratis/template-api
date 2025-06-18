package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type Config map[string]any

// CONFIG STRICTLY STRING
var CONFIG Config = make(map[string]any)

func (c Config) GetString(key string) string {
	value, ok := c[key]
	if !ok {
		log.Fatal().Msgf("Key not found: '%v'", key)
		return ""
	}

	strvalue, ok := value.(string)
	if !ok {
		log.Fatal().Msgf("Failed to assert config with key '%v' as string", key)
		return ""
	}

	return strvalue
}

func (c Config) GetInt(key string) int {
	value, ok := c[key]
	if !ok {
		log.Fatal().Msgf("Key not found: '%v'", key)
		return 0
	}

	result64, ok := value.(float64)
	if !ok {
		log.Fatal().Msgf("Failed to assert value '%v' as number for config key '%v'", value, key)
		return 0
	}

	return int(result64)
}

func (c Config) GetFloat64(key string) float64 {
	value, ok := c[key]
	if !ok {
		log.Fatal().Msgf("Key not found: '%v'", key)
		return 0
	}

	result64, ok := value.(float64)
	if !ok {
		log.Fatal().Msgf("Failed to assert value '%v' as number for config key '%v'", value, key)
		return 0
	}

	return result64
}

func (c Config) SetString(key string, value string) {
	c[key] = value
}

func (c Config) SetFloat64(key string, value float64) {
	c[key] = value
}

func (c Config) SetInt(key string, value int) {
	c.SetFloat64(key, float64(value))
}

func initConfig(ctx context.Context, path string) (err error) {
	cfgcombined := make(map[string]any)

	filePath := path

	log.Info().Msgf("Reading config from: %v", filePath)
	// Read config
	cfg, err := os.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("Failed to read config: %v", err)
	} else {
		// Load config
		err = json.Unmarshal(cfg, &cfgcombined)
		if err != nil {
			return err
		}
		CONFIG = cfgcombined
	}

	// overwrite config with secret
	secret := `
		{"example.secret.format": "kv"}
	`

	// obtain from your external secret provider..
	_ = ctx

	// combine
	err = json.Unmarshal([]byte(secret), &cfgcombined)
	if err != nil {
		log.Error().Msgf("  Failed to parse secret: %v", err)
		return err
	}

	CONFIG = cfgcombined

	return nil
}
