package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

var (
	ErrInvalidListCredentialsRequest = errors.New("invalid list credentials request")
	ErrListCredentials               = errors.New("list credentials")
)

var listCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

type ListCredentialsRequest struct {
	IDs []string `validate:"required,min=1,max=128,dive,required,len=36"`
}

type ListCredentialsResponseCredential struct {
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

type ListCredentialsResponse struct {
	Credentials []*ListCredentialsResponseCredential
}

type ListCredentials interface {
	Exec(ctx context.Context, data *ListCredentialsRequest) (*ListCredentialsResponse, error)
}

type listCredentialsImpl struct {
	dao dao.ListCredentials
}

func (service *listCredentialsImpl) Exec(
	ctx context.Context, data *ListCredentialsRequest,
) (*ListCredentialsResponse, error) {
	var err error

	if err = listCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidListCredentialsRequest, err)
	}

	credentialIDs := make(uuid.UUIDs, len(data.IDs))
	for i, id := range data.IDs {
		credentialIDs[i], err = uuid.Parse(id)
		if err != nil {
			return nil, errors.Join(
				ErrInvalidListCredentialsRequest,
				fmt.Errorf("at position %v: '%s': %w", i, id, err),
			)
		}
	}

	credentials, err := service.dao.Exec(ctx, credentialIDs)
	if err != nil {
		return nil, errors.Join(ErrListCredentials, err)
	}

	response := &ListCredentialsResponse{
		Credentials: lo.Map(credentials, func(item *entities.Credential, _ int) *ListCredentialsResponseCredential {
			return &ListCredentialsResponseCredential{
				ID:                            item.ID.String(),
				Email:                         item.Email,
				Role:                          item.Role,
				EmailValidationTokenID:        item.EmailValidationTokenID,
				PendingEmailValidationTokenID: item.PendingEmailValidationTokenID,
				PasswordTokenID:               item.PasswordTokenID,
				ResetPasswordTokenID:          item.ResetPasswordTokenID,
				CreatedAt:                     item.CreatedAt,
				UpdatedAt:                     item.UpdatedAt,
			}
		}),
	}

	return response, nil
}

func NewListCredentials(dao dao.ListCredentials) ListCredentials {
	return &listCredentialsImpl{dao: dao}
}
