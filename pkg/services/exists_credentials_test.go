package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	daomocks "github.com/a-novel/uservice-credentials/pkg/dao/mocks"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

func TestExistsCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *services.ExistsCredentialsRequest

		shouldCallExistsCredentialsDAO bool
		existsCredentialsDAOResponse   bool
		existsCredentialsDAOError      error

		expect    *services.ExistsCredentialsResponse
		expectErr error
	}{
		{
			name: "OK/ID",

			request: &services.ExistsCredentialsRequest{
				ID: "00000000-0000-0000-0000-000000000001",
			},

			shouldCallExistsCredentialsDAO: true,
			existsCredentialsDAOResponse:   true,

			expect: &services.ExistsCredentialsResponse{
				Exists: true,
			},
		},
		{
			name: "OK/Email",

			request: &services.ExistsCredentialsRequest{
				Email: "user@gmail.com",
			},

			shouldCallExistsCredentialsDAO: true,
			existsCredentialsDAOResponse:   true,

			expect: &services.ExistsCredentialsResponse{
				Exists: true,
			},
		},
		{
			name: "DAO/Error",

			request: &services.ExistsCredentialsRequest{
				Email: "user@gmail.com",
			},

			shouldCallExistsCredentialsDAO: true,
			existsCredentialsDAOError:      errors.New("uwups"),

			expectErr: services.ErrExistsCredentials,
		},
		{
			name: "InvalidRequest/NoData",

			request: &services.ExistsCredentialsRequest{},

			expectErr: services.ErrInvalidExistsCredentialsRequest,
		},
		{
			name: "InvalidRequest/BadID",

			request: &services.ExistsCredentialsRequest{
				ID: "00000000x0000x0000x0000x000000000001",
			},

			expectErr: services.ErrInvalidExistsCredentialsRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			existsCredentialsDAO := daomocks.NewMockExistsCredentials(t)

			if testCase.shouldCallExistsCredentialsDAO {
				existsCredentialsDAO.
					On(
						"Exec",
						context.Background(),
						mock.MatchedBy(func(request *dao.ExistsCredentialsRequest) bool {
							var id uuid.UUID
							if testCase.request.ID != "" {
								id = uuid.MustParse(testCase.request.ID)
							}

							return request.ID == id && request.Email == testCase.request.Email
						}),
					).
					Return(testCase.existsCredentialsDAOResponse, testCase.existsCredentialsDAOError)
			}

			service := services.NewExistsCredentials(existsCredentialsDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			existsCredentialsDAO.AssertExpectations(t)
		})
	}
}
