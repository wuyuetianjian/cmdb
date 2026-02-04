package data

import (
	"context"
	"fmt"

	"cmdb/internal/data/ent"
)

type Data struct {
	Ent *ent.Client
}

func NewEntClient(ctx context.Context, driver, dsn string) (*ent.Client, error) {
	client, err := ent.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open ent client: %w", err)
	}
	if err := client.Schema.Create(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return client, nil
}

