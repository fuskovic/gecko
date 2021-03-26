package main

import (
	"context"
	"net/url"
	"time"

	coder "cdr.dev/coder-cli/coder-sdk"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/xerrors"
)

func main() {
	lambda.Start(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		envVars, err := loadEnvVars(ctx)
		if err != nil {
			return xerrors.Errorf("failed to load environment variables: %w", err)
		}

		var (
			email    = envVars["CODER_EMAIL"]
			password = envVars["CODER_PASSWORD"]
			envName  = envVars["CODER_ENVIRONMENT_NAME"]
		)

		url, err := url.Parse(envVars["CODER_BASE_URL"])
		if err != nil {
			return xerrors.Errorf("failed to parse CODER_BASE_URL: %w", err)
		}

		client, err := coder.NewClient(
			coder.ClientOptions{
				BaseURL:  url,
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			return xerrors.Errorf("failed to initialize coder client: %w", err)
		}

		env, err := getUserEnv(ctx, client, email, envName)
		if err != nil {
			return xerrors.Errorf("failed to get user environment: %w", err)
		}
		return client.RebuildEnvironment(ctx, env.ID)
	})
}
