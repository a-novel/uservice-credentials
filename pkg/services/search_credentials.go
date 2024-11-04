package services

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

var (
	ErrInvalidSearchCredentialsRequest = errors.New("invalid search credentials request")
	ErrSearchCredentials               = errors.New("search credentials")
)

var searchCredentialsValidate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	database.RegisterSortDirection(searchCredentialsValidate)
	entities.RegisterRole(searchCredentialsValidate)
	entities.RegisterSortCredentials(searchCredentialsValidate)
}

type SearchCredentialsRequest struct {
	Limit         int                      `validate:"required,min=1,max=128"`
	Offset        int                      `validate:"omitempty,min=0"`
	Sort          entities.SortCredentials `validate:"omitempty,sort_credentials"`
	SortDirection database.SortDirection   `validate:"omitempty,sort_direction"`
	Emails        []string                 `validate:"omitempty,max=128,dive,email"`
	Roles         []entities.Role          `validate:"omitempty,max=128,dive,role"`
}

type SearchCredentialsResponse struct {
	IDs []string
}

type SearchCredentials interface {
	Exec(ctx context.Context, data *SearchCredentialsRequest) (*SearchCredentialsResponse, error)
}

type searchCredentialsImpl struct {
	dao dao.SearchCredentials
}

func (service *searchCredentialsImpl) Exec(
	ctx context.Context, data *SearchCredentialsRequest,
) (*SearchCredentialsResponse, error) {
	if err := searchCredentialsValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidSearchCredentialsRequest, err)
	}

	ids, err := service.dao.Exec(ctx, &dao.SearchCredentialsRequest{
		Limit:         data.Limit,
		Offset:        data.Offset,
		Sort:          data.Sort,
		SortDirection: data.SortDirection,
		Emails:        data.Emails,
		Roles:         data.Roles,
	})
	if err != nil {
		return nil, errors.Join(ErrSearchCredentials, err)
	}

	return &SearchCredentialsResponse{IDs: ids.Strings()}, nil
}

func NewSearchCredentials(dao dao.SearchCredentials) SearchCredentials {
	return &searchCredentialsImpl{dao: dao}
}
