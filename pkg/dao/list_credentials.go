package dao

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

type ListCredentials interface {
	Exec(ctx context.Context, ids []uuid.UUID) ([]*entities.Credential, error)
}

type listCredentialsImpl struct {
	database bun.IDB
}

func (dao *listCredentialsImpl) Exec(ctx context.Context, ids []uuid.UUID) ([]*entities.Credential, error) {
	credentials := make([]*entities.Credential, 0)

	err := dao.database.NewSelect().Model(&credentials).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return credentials, nil
}

func NewListCredentials(database bun.IDB) ListCredentials {
	return &listCredentialsImpl{database: database}
}
