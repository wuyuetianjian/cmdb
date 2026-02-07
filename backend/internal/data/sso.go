package data

import (
	"context"
	"errors"
	"fmt"

	"cmdb/internal/data/ent"
)

type SSOConfig struct {
	ID       int
	Enabled  bool
	Protocol string
}

type SSOConfigRepo struct {
	ent *ent.Client
}

func NewSSOConfigRepo(entClient *ent.Client) *SSOConfigRepo {
	return &SSOConfigRepo{ent: entClient}
}

func (r *SSOConfigRepo) GetOrCreate(ctx context.Context) (*SSOConfig, error) {
	config, err := r.ent.SSOConfig.Query().Only(ctx)
	if err == nil {
		return mapSSOConfig(config), nil
	}
	if err != nil && !errors.Is(err, ent.ErrNotFound) {
		return nil, fmt.Errorf("query sso config: %w", err)
	}

	config, err = r.ent.SSOConfig.Create().Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create sso config: %w", err)
	}
	return mapSSOConfig(config), nil
}

func (r *SSOConfigRepo) Update(ctx context.Context, enabled bool, protocol string) (*SSOConfig, error) {
	config, err := r.ent.SSOConfig.Query().Only(ctx)
	if err != nil {
		if !errors.Is(err, ent.ErrNotFound) {
			return nil, fmt.Errorf("query sso config: %w", err)
		}
		config, err = r.ent.SSOConfig.Create().Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("create sso config: %w", err)
		}
	}

	if protocol == "" {
		protocol = config.Protocol
	}

	config, err = r.ent.SSOConfig.UpdateOneID(config.ID).
		SetEnabled(enabled).
		SetProtocol(protocol).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update sso config: %w", err)
	}
	return mapSSOConfig(config), nil
}

func mapSSOConfig(config *ent.SSOConfig) *SSOConfig {
	return &SSOConfig{
		ID:       config.ID,
		Enabled:  config.Enabled,
		Protocol: config.Protocol,
	}
}
