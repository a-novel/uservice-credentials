package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

var (
	ErrInvalidUpdateCredentialsRequest = errors.New("invalid update credentials request")
	ErrUpdateCredentials               = errors.New("update credentials")
)

var updateCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	entities.RegisterRole(updateCredentialsValidate)
}

type UpdateCredentialsRequest struct {
	ID                            string        `validate:"required,len=36"`
	Email                         string        `validate:"required,email,max=256"`
	Role                          entities.Role `validate:"omitempty,role"`
	EmailValidationTokenID        string        `validate:"omitempty,min=1,max=128"`
	PendingEmailValidationTokenID string        `validate:"omitempty,min=1,max=128"`
	PasswordTokenID               string        `validate:"omitempty,min=1,max=128"`
	ResetPasswordTokenID          string        `validate:"omitempty,min=1,max=128"`
}

type UpdateCredentialsResponse struct {
	ID    string
	Email string
	Role  entities.Role

	EmailValidationTokenID        string
	PendingEmailValidationTokenID string
	PasswordTokenID               string
	ResetPasswordTokenID          string

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UpdateCredentials interface {
	Exec(ctx context.Context, data *UpdateCredentialsRequest) (*UpdateCredentialsResponse, error)
}

type updateCredentialsImpl struct {
	dao dao.UpdateCredentials
}

func (service *updateCredentialsImpl) Exec(
	ctx context.Context, data *UpdateCredentialsRequest,
) (*UpdateCredentialsResponse, error) {
	if err := updateCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidUpdateCredentialsRequest, err)
	}

	credentialsID, err := uuid.Parse(data.ID)
	if err != nil {
		return nil, errors.Join(ErrInvalidUpdateCredentialsRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
	}

	credentials, err := service.dao.Exec(ctx, credentialsID, time.Now(), &dao.UpdateCredentialsRequest{
		Email:                         data.Email,
		Role:                          data.Role,
		EmailValidationTokenID:        data.EmailValidationTokenID,
		PendingEmailValidationTokenID: data.PendingEmailValidationTokenID,
		PasswordTokenID:               data.PasswordTokenID,
		ResetPasswordTokenID:          data.ResetPasswordTokenID,
	})
	if err != nil {
		return nil, errors.Join(ErrUpdateCredentials, err)
	}

	return &UpdateCredentialsResponse{
		ID:                            credentials.ID.String(),
		Email:                         credentials.Email,
		Role:                          credentials.Role,
		EmailValidationTokenID:        credentials.EmailValidationTokenID,
		PendingEmailValidationTokenID: credentials.PendingEmailValidationTokenID,
		PasswordTokenID:               credentials.PasswordTokenID,
		ResetPasswordTokenID:          credentials.ResetPasswordTokenID,
		CreatedAt:                     credentials.CreatedAt,
		UpdatedAt:                     credentials.UpdatedAt,
	}, nil
}

func NewUpdateCredentials(dao dao.UpdateCredentials) UpdateCredentials {
	return &updateCredentialsImpl{dao: dao}
}
