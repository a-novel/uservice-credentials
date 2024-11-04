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
	ErrInvalidGetCredentialsRequest = errors.New("invalid get credentials request")
	ErrGetCredentials               = errors.New("get credentials")
)

var getCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

type GetCredentialsRequest struct {
	ID    string `validate:"required_without=Email,omitempty,len=36"`
	Email string `validate:"required_without=ID,omitempty,email,max=256"`
}

type GetCredentialsResponse struct {
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

type GetCredentials interface {
	Exec(ctx context.Context, data *GetCredentialsRequest) (*GetCredentialsResponse, error)
}

type getCredentialsImpl struct {
	dao dao.GetCredentials
}

func (service *getCredentialsImpl) Exec(
	ctx context.Context, data *GetCredentialsRequest,
) (*GetCredentialsResponse, error) {
	var err error

	if err = getCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidGetCredentialsRequest, err)
	}

	var credentialsID uuid.UUID
	if data.ID != "" {
		credentialsID, err = uuid.Parse(data.ID)
		if err != nil {
			return nil, errors.Join(ErrInvalidGetCredentialsRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
		}
	}

	request := &dao.GetCredentialsRequest{
		Email: data.Email,
		ID:    credentialsID,
	}

	credentials, err := service.dao.Exec(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrGetCredentials, err)
	}

	return &GetCredentialsResponse{
		ID:    credentials.ID.String(),
		Email: credentials.Email,
		Role:  credentials.Role,

		EmailValidationTokenID:        credentials.EmailValidationTokenID,
		PendingEmailValidationTokenID: credentials.PendingEmailValidationTokenID,
		PasswordTokenID:               credentials.PasswordTokenID,
		ResetPasswordTokenID:          credentials.ResetPasswordTokenID,

		CreatedAt: credentials.CreatedAt,
		UpdatedAt: credentials.UpdatedAt,
	}, nil
}

func NewGetCredentials(dao dao.GetCredentials) GetCredentials {
	return &getCredentialsImpl{dao: dao}
}
