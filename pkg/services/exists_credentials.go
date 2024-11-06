package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/a-novel/uservice-credentials/pkg/dao"
)

var (
	ErrInvalidExistsCredentialsRequest = errors.New("invalid exists credentials request")
	ErrExistsCredentials               = errors.New("exists credentials")
)

var existsCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

type ExistsCredentialsRequest struct {
	ID    string `validate:"required_without=Email,omitempty,len=36"`
	Email string `validate:"required_without=ID,omitempty,email,max=256"`
}

type ExistsCredentialsResponse struct {
	Exists bool
}

type ExistsCredentials interface {
	Exec(ctx context.Context, data *ExistsCredentialsRequest) (*ExistsCredentialsResponse, error)
}

type existsCredentialsImpl struct {
	dao dao.ExistsCredentials
}

func (service *existsCredentialsImpl) Exec(
	ctx context.Context, data *ExistsCredentialsRequest,
) (*ExistsCredentialsResponse, error) {
	var err error

	if err = existsCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidExistsCredentialsRequest, err)
	}

	var credentialsID uuid.UUID
	if data.ID != "" {
		credentialsID, err = uuid.Parse(data.ID)
		if err != nil {
			return nil, errors.Join(ErrInvalidExistsCredentialsRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
		}
	}

	request := &dao.ExistsCredentialsRequest{
		Email: data.Email,
		ID:    credentialsID,
	}

	exists, err := service.dao.Exec(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrExistsCredentials, err)
	}

	return &ExistsCredentialsResponse{Exists: exists}, nil
}

func NewExistsCredentials(dao dao.ExistsCredentials) ExistsCredentials {
	return &existsCredentialsImpl{dao: dao}
}
