package core

import (
	"fmt"
	"os"

	envmanager_dtos "github.com/keyloom/web-api/dtos/env-manager"
)

type EnvManager struct{}

func (e *EnvManager) ValidateEnvs(vars []string) ([]string, error) {
	var values []string
	var missingVars []string
	for _, v := range vars {
		envValue := os.Getenv(v)
		if envValue == "" {
			missingVars = append(missingVars, v)
		} else {
			values = append(values, envValue)
		}
	}
	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing environment variables: %v", missingVars)
	}
	return values, nil
}

func (e *EnvManager) ValidateEnv(name string) (string, error) {
	envValue := os.Getenv(name)
	if envValue == "" {
		return "", fmt.Errorf("missing environment variable: %s", name)
	}
	return envValue, nil
}

func (e *EnvManager) GetMongoConfig() (envmanager_dtos.MongoConfig, error) {
	vars := []string{
		"MONGODB_HOST",
		"MONGODB_PORT",
		"MONGODB_DB",
		"MONGODB_USER",
		"MONGODB_PASSWORD",
		"MONGODB_AUTH_SOURCE",
	}
	values, err := e.ValidateEnvs(vars)
	if err != nil {
		return envmanager_dtos.MongoConfig{}, err
	}
	mongoConfig := envmanager_dtos.MongoConfig{
		Host:         values[0],
		Port:         values[1],
		DatabaseName: values[2],
		Username:     values[3],
		Password:     values[4],
		AuthSource:   values[5],
	}
	return mongoConfig, nil
}
