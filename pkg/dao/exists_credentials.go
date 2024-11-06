package dao

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

type ExistsCredentialsRequest struct {
	Email string
	ID    uuid.UUID
}

type ExistsCredentials interface {
	Exec(ctx context.Context, request *ExistsCredentialsRequest) (bool, error)
}

type existsCredentialsImpl struct {
	database bun.IDB
}

func (dao *existsCredentialsImpl) Exec(ctx context.Context, request *ExistsCredentialsRequest) (bool, error) {
	query := dao.database.NewSelect().Model((*entities.Credential)(nil))

	if request.Email == "" && request.ID == uuid.Nil {
		return false, ErrCredentialsNotFound
	}

	if request.Email != "" {
		query.Where("email = ?", request.Email)
	}
	if request.ID != uuid.Nil {
		query.Where("id = ?", request.ID)
	}

	ok, err := query.Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("exec query: %w", err)
	}

	return ok, nil
}

func NewExistsCredentials(database bun.IDB) ExistsCredentials {
	return &existsCredentialsImpl{database: database}
}
