package main

import (
	"context"
	"os"

	"golang.org/x/xerrors"
)

var requiredEnvVars = []string{
	"CODER_EMAIL",
	"CODER_PASSWORD",
	"CODER_BASE_URL",
	"CODER_ENVIRONMENT_NAME",
}

func loadEnvVars(_ context.Context) (map[string]string, error) {
	envVars := make(map[string]string)
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			return nil, xerrors.Errorf("required environment variable %q is unset", envVar)
		}
		envVars[envVar] = value
	}
	return envVars, nil
}
