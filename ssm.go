package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSMStore implements Store for AWS SSM Parameter Store.
type SSMStore struct {
	client *ssm.Client
}

func NewSSMStore(cfg aws.Config) *SSMStore {
	return &SSMStore{client: ssm.NewFromConfig(cfg)}
}

func (s *SSMStore) GetAll(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)
	var nextToken *string

	for {
		out, err := s.client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String("/"),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			NextToken:      nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("ssm GetParametersByPath: %w", err)
		}
		for _, p := range out.Parameters {
			result[*p.Name] = *p.Value
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return result, nil
}

func (s *SSMStore) Put(ctx context.Context, key, value string) error {
	_, err := s.client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("ssm PutParameter %s: %w", key, err)
	}
	return nil
}

func (s *SSMStore) Delete(ctx context.Context, keys []string) error {
	// Batch delete max 10 per request
	for i := 0; i < len(keys); i += 10 {
		end := i + 10
		if end > len(keys) {
			end = len(keys)
		}
		_, err := s.client.DeleteParameters(ctx, &ssm.DeleteParametersInput{
			Names: keys[i:end],
		})
		if err != nil {
			return fmt.Errorf("ssm DeleteParameters: %w", err)
		}
	}
	return nil
}
