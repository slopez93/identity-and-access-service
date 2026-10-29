package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	err := env.ParseWithOptions(
		&cfg,
		env.Options{
			FuncMap: map[reflect.Type]env.ParserFunc{
				reflect.TypeOf(CognitoServiceClientsConfig{}): parseServiceClients,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	return &cfg, nil
}

func parseServiceClients(
	value string,
) (interface{}, error) {
	if value == "" {
		return CognitoServiceClientsConfig{}, nil
	}

	var clients CognitoServiceClientsConfig

	if err := json.Unmarshal(
		[]byte(value),
		&clients,
	); err != nil {
		return nil, fmt.Errorf(
			"parse COGNITO_SERVICE_CLIENTS: %w",
			err,
		)
	}

	return clients, nil
}
