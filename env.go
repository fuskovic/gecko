package main

import (
	"context"
	"fmt"
	"os"
)

var requiredEnvVars = []string{
	"CODER_EMAIL",
	"CODER_PASSWORD",
	"CODER_ACCESS_URL",
	"CODER_ENVIRONMENT_NAME",
}

func loadEnvVars(_ context.Context) (map[string]string, error) {
	envVars := make(map[string]string)
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			return nil, fmt.Errorf("required environment variable %q is unset", envVar)
		}
		envVars[envVar] = value
	}
	return envVars, nil
}
