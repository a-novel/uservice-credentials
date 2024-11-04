package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

type GetCredentialsRequest struct {
	Email string
	ID    uuid.UUID
}

type GetCredentials interface {
	Exec(ctx context.Context, request *GetCredentialsRequest) (*entities.Credential, error)
}

type getCredentialsImpl struct {
	database bun.IDB
}

func (dao *getCredentialsImpl) Exec(ctx context.Context, request *GetCredentialsRequest) (*entities.Credential, error) {
	credential := new(entities.Credential)

	query := dao.database.NewSelect().Model(credential)

	if request.Email == "" && request.ID == uuid.Nil {
		return nil, ErrCredentialsNotFound
	}

	if request.Email != "" {
		query.Where("email = ?", request.Email)
	}
	if request.ID != uuid.Nil {
		query.Where("id = ?", request.ID)
	}

	err := query.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCredentialsNotFound
		}

		return nil, fmt.Errorf("exec query: %w", err)
	}

	return credential, nil
}

func NewGetCredentials(database bun.IDB) GetCredentials {
	return &getCredentialsImpl{database: database}
}
