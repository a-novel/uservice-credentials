package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	daomocks "github.com/a-novel/uservice-credentials/pkg/dao/mocks"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

func TestGetCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *services.GetCredentialsRequest

		shouldCallGetCredentialsDAO bool
		getCredentialsDAOResponse   *entities.Credential
		getCredentialsDAOError      error

		expect    *services.GetCredentialsResponse
		expectErr error
	}{
		{
			name: "OK/ID",

			request: &services.GetCredentialsRequest{
				ID: "00000000-0000-0000-0000-000000000001",
			},

			shouldCallGetCredentialsDAO: true,
			getCredentialsDAOResponse: &entities.Credential{
				ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Email:                         "user@gmail.com",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenID: "00000000-0000-0000-0000-000000000005",
				PasswordTokenID:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.GetCredentialsResponse{
				ID:                            "00000000-0000-0000-0000-000000000004",
				Email:                         "user@gmail.com",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenID: "00000000-0000-0000-0000-000000000005",
				PasswordTokenID:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "OK/Email",

			request: &services.GetCredentialsRequest{
				Email: "user@gmail.com",
			},

			shouldCallGetCredentialsDAO: true,
			getCredentialsDAOResponse: &entities.Credential{
				ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Email:                         "user@gmail.com",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenID: "00000000-0000-0000-0000-000000000005",
				PasswordTokenID:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.GetCredentialsResponse{
				ID:                            "00000000-0000-0000-0000-000000000004",
				Email:                         "user@gmail.com",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenID: "00000000-0000-0000-0000-000000000005",
				PasswordTokenID:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "DAO/Error",

			request: &services.GetCredentialsRequest{
				Email: "user@gmail.com",
			},

			shouldCallGetCredentialsDAO: true,
			getCredentialsDAOError:      errors.New("uwups"),

			expectErr: services.ErrGetCredentials,
		},
		{
			name: "InvalidRequest/NoData",

			request: &services.GetCredentialsRequest{},

			expectErr: services.ErrInvalidGetCredentialsRequest,
		},
		{
			name: "InvalidRequest/BadID",

			request: &services.GetCredentialsRequest{
				ID: "00000000x0000x0000x0000x000000000001",
			},

			expectErr: services.ErrInvalidGetCredentialsRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			getCredentialsDAO := daomocks.NewMockGetCredentials(t)

			if testCase.shouldCallGetCredentialsDAO {
				getCredentialsDAO.
					On(
						"Exec",
						context.Background(),
						mock.MatchedBy(func(request *dao.GetCredentialsRequest) bool {
							var id uuid.UUID
							if testCase.request.ID != "" {
								id = uuid.MustParse(testCase.request.ID)
							}

							return request.ID == id && request.Email == testCase.request.Email
						}),
					).
					Return(testCase.getCredentialsDAOResponse, testCase.getCredentialsDAOError)
			}

			service := services.NewGetCredentials(getCredentialsDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			getCredentialsDAO.AssertExpectations(t)
		})
	}
}
