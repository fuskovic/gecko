package main

import (
	"context"
	"fmt"

	coder "cdr.dev/coder-cli/coder-sdk"
)

func getUserEnv(ctx context.Context, client coder.Client, email, envName string) (*coder.Environment, error) {
	envs, err := listUserEnvs(ctx, client, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments for %q: %w", email, err)
	}

	var env *coder.Environment
	for _, e := range envs {
		if e.Name == envName {
			env = &e
			break
		}
	}

	if env == nil {
		return nil, fmt.Errorf("environment %q not found", envName)
	}
	return env, nil
}

func listUserEnvs(ctx context.Context, client coder.Client, email string) ([]coder.Environment, error) {
	user, err := client.UserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email %q: %w", email, err)
	}

	allOrgs, err := client.Organizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	var (
		userOrgs = listUserOrgs(user, allOrgs)
		userEnvs = []coder.Environment{}
	)

	for _, org := range userOrgs {
		envs, err := client.UserEnvironmentsByOrganization(ctx, user.ID, org.ID)
		if err != nil {
			return nil, fmt.Errorf("faield to get environments for %q: %w", org.Name, err)
		}
		userEnvs = append(userEnvs, envs...)
	}
	return userEnvs, nil
}

func listUserOrgs(user *coder.User, orgs []coder.Organization) []coder.Organization {
	var userOrgs []coder.Organization
	for _, org := range orgs {
		for _, member := range org.Members {
			if member.ID == user.ID {
				userOrgs = append(userOrgs, org)
				break
			}
		}
	}
	return userOrgs
}
