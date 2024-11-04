package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

type CreateCredentialsRequest struct {
	Email                  string
	Role                   entities.Role
	EmailValidationTokenID string
	PasswordTokenID        string
	ResetPasswordTokenID   string
}

type CreateCredentials interface {
	Exec(
		ctx context.Context, id uuid.UUID, now time.Time, request *CreateCredentialsRequest,
	) (*entities.Credential, error)
}

type createCredentialsImpl struct {
	database bun.IDB
}

func (dao *createCredentialsImpl) Exec(
	ctx context.Context, id uuid.UUID, now time.Time, request *CreateCredentialsRequest,
) (*entities.Credential, error) {
	model := &entities.Credential{
		ID:                     id,
		Email:                  request.Email,
		Role:                   request.Role,
		EmailValidationTokenID: request.EmailValidationTokenID,
		PasswordTokenID:        request.PasswordTokenID,
		ResetPasswordTokenID:   request.ResetPasswordTokenID,
		CreatedAt:              now,
	}

	_, err := dao.database.NewInsert().Model(model).Returning("*").Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return nil, ErrCredentialsAlreadyExist
		}

		return nil, fmt.Errorf("exec query: %w", err)
	}

	return model, nil
}

func NewCreateCredentials(database bun.IDB) CreateCredentials {
	return &createCredentialsImpl{database: database}
}
