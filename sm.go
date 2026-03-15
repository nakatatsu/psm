package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smtypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// SMStore implements Store for AWS Secrets Manager.
type SMStore struct {
	client *secretsmanager.Client
}

func NewSMStore(cfg aws.Config) *SMStore {
	return &SMStore{client: secretsmanager.NewFromConfig(cfg)}
}

func (s *SMStore) GetAll(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)

	// List all secret names
	var names []string
	var nextToken *string
	for {
		out, err := s.client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("sm ListSecrets: %w", err)
		}
		for _, s := range out.SecretList {
			names = append(names, *s.Name)
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	// Get each secret value individually (avoids BatchGetSecretValue permission)
	for _, name := range names {
		out, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return nil, fmt.Errorf("sm GetSecretValue %s: %w", name, err)
		}
		if out.SecretString != nil {
			result[*out.Name] = *out.SecretString
		}
	}

	return result, nil
}

func (s *SMStore) Put(ctx context.Context, key, value string) error {
	_, err := s.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(key),
		SecretString: aws.String(value),
	})
	if err != nil {
		var ree *smtypes.ResourceExistsException
		if errors.As(err, &ree) {
			_, err = s.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
				SecretId:     aws.String(key),
				SecretString: aws.String(value),
			})
			if err != nil {
				return fmt.Errorf("sm PutSecretValue %s: %w", key, err)
			}
			return nil
		}
		return fmt.Errorf("sm CreateSecret %s: %w", key, err)
	}
	return nil
}

func (s *SMStore) Delete(ctx context.Context, keys []string) error {
	for _, key := range keys {
		_, err := s.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(key),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("sm DeleteSecret %s: %w", key, err)
		}
	}
	return nil
}
