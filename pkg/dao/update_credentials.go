package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

type UpdateCredentialsRequest struct {
	Email string
	Role  entities.Role

	EmailValidationTokenID        string
	PendingEmailValidationTokenID string
	PasswordTokenID               string
	ResetPasswordTokenID          string
}

type UpdateCredentials interface {
	Exec(
		ctx context.Context, id uuid.UUID, now time.Time, data *UpdateCredentialsRequest,
	) (*entities.Credential, error)
}

type updateCredentialsImpl struct {
	database bun.IDB
}

func (dao *updateCredentialsImpl) Exec(
	ctx context.Context, id uuid.UUID, now time.Time, data *UpdateCredentialsRequest,
) (*entities.Credential, error) {
	model := &entities.Credential{
		ID:                            id,
		Email:                         data.Email,
		Role:                          data.Role,
		EmailValidationTokenID:        data.EmailValidationTokenID,
		PendingEmailValidationTokenID: data.PendingEmailValidationTokenID,
		PasswordTokenID:               data.PasswordTokenID,
		ResetPasswordTokenID:          data.ResetPasswordTokenID,
		UpdatedAt:                     &now,
	}

	res, err := dao.database.
		NewUpdate().
		Model(model).
		WherePK().
		ExcludeColumn("id", "created_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return nil, ErrCredentialsNotFound
	}

	return model, nil
}

func NewUpdateCredentials(database bun.IDB) UpdateCredentials {
	return &updateCredentialsImpl{database: database}
}
