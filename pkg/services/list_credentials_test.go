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

	daomocks "github.com/a-novel/uservice-credentials/pkg/dao/mocks"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

func TestListCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *services.ListCredentialsRequest

		shouldCallListCredentialsDAO bool
		listCredentialsDAOResponse   []*entities.Credential
		listCredentialsDAOError      error

		expect    *services.ListCredentialsResponse
		expectErr error
	}{
		{
			name: "OK",

			request: &services.ListCredentialsRequest{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
				},
			},

			shouldCallListCredentialsDAO: true,
			listCredentialsDAOResponse: []*entities.Credential{
				{
					ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Email:                         "email-1",
					Role:                          entities.RoleCore,
					EmailValidationTokenID:        "email-validation-token-id-1",
					PendingEmailValidationTokenID: "pending-email-validation-token-id-1",
					PasswordTokenID:               "password-token-id-1",
					ResetPasswordTokenID:          "reset-password-token-id-1",
					CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Email:     "email-3",
					Role:      entities.RoleNone,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &services.ListCredentialsResponse{
				Credentials: []*services.ListCredentialsResponseCredential{
					{
						ID:                            "00000000-0000-0000-0000-000000000001",
						Email:                         "email-1",
						Role:                          entities.RoleCore,
						EmailValidationTokenID:        "email-validation-token-id-1",
						PendingEmailValidationTokenID: "pending-email-validation-token-id-1",
						PasswordTokenID:               "password-token-id-1",
						ResetPasswordTokenID:          "reset-password-token-id-1",
						CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        "00000000-0000-0000-0000-000000000003",
						Email:     "email-3",
						Role:      entities.RoleNone,
						CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "OK/NoReturn",

			request: &services.ListCredentialsRequest{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
				},
			},

			shouldCallListCredentialsDAO: true,
			listCredentialsDAOResponse:   []*entities.Credential{},

			expect: &services.ListCredentialsResponse{
				Credentials: []*services.ListCredentialsResponseCredential{},
			},
		},
		{
			name: "DAO/Error",

			request: &services.ListCredentialsRequest{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
				},
			},

			shouldCallListCredentialsDAO: true,
			listCredentialsDAOError:      errors.New("uwups"),

			expectErr: services.ErrListCredentials,
		},
		{
			name: "InvalidRequest/NoID",

			request: &services.ListCredentialsRequest{},

			expectErr: services.ErrInvalidListCredentialsRequest,
		},
		{
			name: "InvalidRequest/InvalidID",

			request: &services.ListCredentialsRequest{
				IDs: []string{
					"00000000-0000-0000-0000-000000000001",
					"00000000x0000x0000x0000x000000000002",
				},
			},

			expectErr: services.ErrInvalidListCredentialsRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			listCredentialsDAO := daomocks.NewMockListCredentials(t)

			if testCase.shouldCallListCredentialsDAO {
				listCredentialsDAO.
					On(
						"Exec",
						context.Background(),
						mock.MatchedBy(func(ids uuid.UUIDs) bool {
							strIDs := ids.Strings()
							if len(strIDs) != len(testCase.request.IDs) {
								return false
							}

							for index, id := range testCase.request.IDs {
								if id != strIDs[index] {
									return false
								}
							}

							return true
						}),
					).
					Return(testCase.listCredentialsDAOResponse, testCase.listCredentialsDAOError)
			}

			service := services.NewListCredentials(listCredentialsDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			listCredentialsDAO.AssertExpectations(t)
		})
	}
}
