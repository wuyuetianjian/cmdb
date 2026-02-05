package ent

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("ent: not found")

type Client struct {
	Schema     *Schema
	SSOConfig  *SSOConfigClient
}

type Schema struct{}

func (s *Schema) Create(_ context.Context, _ ...interface{}) error {
	return nil
}

func Open(_ string, _ string) (*Client, error) {
	client := &Client{
		Schema:    &Schema{},
		SSOConfig: NewSSOConfigClient(),
	}
	return client, nil
}

func (c *Client) Close() error {
	return nil
}
