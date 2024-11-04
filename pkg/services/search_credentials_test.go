package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	daomocks "github.com/a-novel/uservice-credentials/pkg/dao/mocks"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

func TestSearchCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *services.SearchCredentialsRequest

		shouldCallSearchCredentialsDAO bool
		searchCredentialsDAOResponse   uuid.UUIDs
		searchCredentialsDAOError      error

		expect    *services.SearchCredentialsResponse
		expectErr error
	}{
		{
			name: "OK/All",

			request: &services.SearchCredentialsRequest{
				Limit:         10,
				Offset:        0,
				Sort:          entities.SortCredentialsEmail,
				SortDirection: database.SortDirectionAsc,
				Emails:        []string{"email-1@gmail.com", "email-2@gmail.com"},
				Roles:         []entities.Role{entities.RoleCore, entities.RoleAdmin},
			},

			shouldCallSearchCredentialsDAO: true,
			searchCredentialsDAOResponse: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},

			expect: &services.SearchCredentialsResponse{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
				},
			},
		},
		{
			name: "OK/Minimal",

			request: &services.SearchCredentialsRequest{
				Limit: 10,
			},

			shouldCallSearchCredentialsDAO: true,
			searchCredentialsDAOResponse: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},

			expect: &services.SearchCredentialsResponse{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
				},
			},
		},

		{
			name: "DAO/Error",

			request: &services.SearchCredentialsRequest{
				Limit: 10,
			},

			shouldCallSearchCredentialsDAO: true,
			searchCredentialsDAOError:      errors.New("uwups"),

			expectErr: services.ErrSearchCredentials,
		},
		{
			name: "InvalidRequest/InvalidSort",

			request: &services.SearchCredentialsRequest{
				Limit: 10,
				Sort:  "invalid",
			},

			expectErr: services.ErrInvalidSearchCredentialsRequest,
		},
		{
			name: "InvalidRequest/InvalidSortDirection",

			request: &services.SearchCredentialsRequest{
				Limit:         10,
				SortDirection: "invalid",
			},

			expectErr: services.ErrInvalidSearchCredentialsRequest,
		},
		{
			name: "InvalidRequest/LimitTooLow",

			request: &services.SearchCredentialsRequest{
				Limit: 0,
			},

			expectErr: services.ErrInvalidSearchCredentialsRequest,
		},
		{
			name: "InvalidRequest/LimitTooHigh",

			request: &services.SearchCredentialsRequest{
				Limit: 129,
			},

			expectErr: services.ErrInvalidSearchCredentialsRequest,
		},
		{
			name: "InvalidRequest/OffsetTooLow",

			request: &services.SearchCredentialsRequest{
				Limit:  10,
				Offset: -1,
			},

			expectErr: services.ErrInvalidSearchCredentialsRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			searchCredentialsDAO := daomocks.NewMockSearchCredentials(t)

			if testCase.shouldCallSearchCredentialsDAO {
				searchCredentialsDAO.
					On("Exec", context.Background(), &dao.SearchCredentialsRequest{
						Limit:         testCase.request.Limit,
						Offset:        testCase.request.Offset,
						Sort:          testCase.request.Sort,
						SortDirection: testCase.request.SortDirection,
						Emails:        testCase.request.Emails,
						Roles:         testCase.request.Roles,
					}).
					Return(testCase.searchCredentialsDAOResponse, testCase.searchCredentialsDAOError)
			}

			service := services.NewSearchCredentials(searchCredentialsDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			searchCredentialsDAO.AssertExpectations(t)
		})
	}
}
