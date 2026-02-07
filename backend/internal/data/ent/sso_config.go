package ent

import (
	"context"
)

type SSOConfig struct {
	ID       int
	Enabled  bool
	Protocol string
}

type SSOConfigClient struct {
	config *SSOConfig
}

func NewSSOConfigClient() *SSOConfigClient {
	return &SSOConfigClient{
		config: &SSOConfig{
			ID:       1,
			Enabled:  false,
			Protocol: "saml2",
		},
	}
}

func (c *SSOConfigClient) Query() *SSOConfigQuery {
	return &SSOConfigQuery{client: c}
}

func (c *SSOConfigClient) Create() *SSOConfigCreate {
	return &SSOConfigCreate{client: c}
}

func (c *SSOConfigClient) UpdateOneID(id int) *SSOConfigUpdateOne {
	return &SSOConfigUpdateOne{client: c, id: id}
}

type SSOConfigQuery struct {
	client *SSOConfigClient
}

func (q *SSOConfigQuery) Only(_ context.Context) (*SSOConfig, error) {
	if q.client.config == nil {
		return nil, ErrNotFound
	}
	return q.client.config, nil
}

type SSOConfigCreate struct {
	client *SSOConfigClient
}

func (c *SSOConfigCreate) Save(_ context.Context) (*SSOConfig, error) {
	if c.client.config == nil {
		c.client.config = &SSOConfig{
			ID:       1,
			Enabled:  false,
			Protocol: "saml2",
		}
	}
	return c.client.config, nil
}

type SSOConfigUpdateOne struct {
	client   *SSOConfigClient
	id       int
	enabled  *bool
	protocol *string
}

func (u *SSOConfigUpdateOne) SetEnabled(value bool) *SSOConfigUpdateOne {
	u.enabled = &value
	return u
}

func (u *SSOConfigUpdateOne) SetProtocol(value string) *SSOConfigUpdateOne {
	u.protocol = &value
	return u
}

func (u *SSOConfigUpdateOne) Save(_ context.Context) (*SSOConfig, error) {
	if u.client.config == nil {
		return nil, ErrNotFound
	}
	if u.client.config.ID != u.id {
		return nil, ErrNotFound
	}
	if u.enabled != nil {
		u.client.config.Enabled = *u.enabled
	}
	if u.protocol != nil {
		u.client.config.Protocol = *u.protocol
	}
	return u.client.config, nil
}
