package services

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

var (
	ErrInvalidCreateCredentialsRequest = errors.New("invalid create credentials request")
	ErrCreateCredentials               = errors.New("create credentials")
)

var createCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	entities.RegisterRole(createCredentialsValidate)
}

type CreateCredentialsRequest struct {
	Email                  string        `validate:"required,email,max=256"`
	Role                   entities.Role `validate:"omitempty,role"`
	EmailValidationTokenID string        `validate:"omitempty,min=1,max=128"`
	PasswordTokenID        string        `validate:"omitempty,min=1,max=128"`
	ResetPasswordTokenID   string        `validate:"omitempty,min=1,max=128"`
}

type CreateCredentialsResponse struct {
	ID    string
	Email string
	Role  entities.Role

	EmailValidationTokenID string
	PasswordTokenID        string
	ResetPasswordTokenID   string

	CreatedAt time.Time
}

type CreateCredentials interface {
	Exec(ctx context.Context, data *CreateCredentialsRequest) (*CreateCredentialsResponse, error)
}

type createCredentialsImpl struct {
	dao dao.CreateCredentials
}

func (service *createCredentialsImpl) Exec(
	ctx context.Context, data *CreateCredentialsRequest,
) (*CreateCredentialsResponse, error) {
	if err := createCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidCreateCredentialsRequest, err)
	}

	request := &dao.CreateCredentialsRequest{
		Email:                  data.Email,
		Role:                   data.Role,
		EmailValidationTokenID: data.EmailValidationTokenID,
		PasswordTokenID:        data.PasswordTokenID,
		ResetPasswordTokenID:   data.ResetPasswordTokenID,
	}

	res, err := service.dao.Exec(ctx, uuid.New(), time.Now(), request)
	if err != nil {
		return nil, errors.Join(ErrCreateCredentials, err)
	}

	return &CreateCredentialsResponse{
		ID:    res.ID.String(),
		Email: res.Email,
		Role:  res.Role,

		EmailValidationTokenID: res.EmailValidationTokenID,
		PasswordTokenID:        res.PasswordTokenID,
		ResetPasswordTokenID:   res.ResetPasswordTokenID,

		CreatedAt: res.CreatedAt,
	}, nil
}

func NewCreateCredentials(dao dao.CreateCredentials) CreateCredentials {
	return &createCredentialsImpl{dao: dao}
}
