package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	coder "cdr.dev/coder-cli/coder-sdk"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/xerrors"
)

func main() {
	lambda.Start(func(ctx context.Context) (AlexaResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		resp := newResponse()

		envVars, err := loadEnvVars(ctx)
		if err != nil {
			resp.Say("hmmm, it looks like I couldn't load your environment variables")
			return resp, xerrors.Errorf("failed to load environment variables: %w", err)
		}

		var (
			email    = envVars["CODER_EMAIL"]
			password = envVars["CODER_PASSWORD"]
			envName  = envVars["CODER_ENVIRONMENT_NAME"]
		)

		url, err := url.Parse(envVars["CODER_ACCESS_URL"])
		if err != nil {
			resp.Say("hmmm, it looks I couldn't parse the coder access url you set")
			return resp, xerrors.Errorf("failed to parse CODER_ACCESS_URL: %w", err)
		}

		client, err := coder.NewClient(
			coder.ClientOptions{
				BaseURL:  url,
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			resp.Say("hmmm, it looks I couldn't initialize the coder client")
			return resp, xerrors.Errorf("failed to initialize coder client: %w", err)
		}

		env, err := getUserEnv(ctx, client, email, envName)
		if err != nil {
			resp.Say(fmt.Sprintf("hmmm, it looks I couldn't find the %s environment", envName))
			return resp, xerrors.Errorf("failed to get user environment: %w", err)
		}

		if err := client.RebuildEnvironment(ctx, env.ID); err != nil {
			resp.Say("hmmm, something wrong when I tried adding an environment build job to the queue")
			return resp, xerrors.Errorf("failed to enqueue environment build job: %w", err)
		}
		resp.Say(fmt.Sprintf("OK, I added a new environment build job to the queue for the %s environment", envName))
		return resp, nil
	})
}
