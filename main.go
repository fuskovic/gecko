package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	coder "cdr.dev/coder-cli/coder-sdk"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context) (*AlexaResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		envVars, err := loadEnvVars(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load environment variables: %w", err)
		}

		var (
			email    = envVars["CODER_EMAIL"]
			password = envVars["CODER_PASSWORD"]
			envName  = envVars["CODER_ENVIRONMENT_NAME"]
		)

		url, err := url.Parse(envVars["CODER_ACCESS_URL"])
		if err != nil {
			return nil, fmt.Errorf("failed to parse CODER_ACCESS_URL: %w", err)
		}

		client, err := coder.NewClient(
			coder.ClientOptions{
				BaseURL:  url,
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("failed to initialize coder client: %w", err)
		}

		env, err := getUserEnv(ctx, client, email, envName)
		if err != nil {
			return nil, fmt.Errorf("failed to get user environment: %w", err)
		}

		if err := client.RebuildEnvironment(ctx, env.ID); err != nil {
			return nil, fmt.Errorf("failed to enqueue environment build job: %w", err)
		}

		return newResponse(
			fmt.Sprintf("OK, I added a new environment build job to the queue for the %s environment", envName),
		), nil
	})
}
