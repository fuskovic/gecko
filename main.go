package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	coder "cdr.dev/coder-cli/coder-sdk"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/xerrors"
)

func main() {
	lambda.Start(func(ctx context.Context, req AlexaRequest) (AlexaResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		log.Printf("checking intent: %s", req.Intent.Name)
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

		url, err := url.Parse(envVars["CODER_BASE_URL"])
		if err != nil {
			resp.Say("hmmm, it looks I couldn't parse the coder base url you set")
			return resp, xerrors.Errorf("failed to parse CODER_BASE_URL: %w", err)
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
		resp.Say(fmt.Sprintf("OK, i'm rebuilding %s", envName))
		return resp, client.RebuildEnvironment(ctx, env.ID)
	})
}
