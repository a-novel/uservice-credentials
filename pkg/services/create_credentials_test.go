package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	daomocks "github.com/a-novel/uservice-credentials/pkg/dao/mocks"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

func TestCreateCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *services.CreateCredentialsRequest

		shouldCallCreateCredentialsDAO bool
		createCredentialsDAOResponse   *entities.Credential
		createCredentialsDAOError      error

		expect    *services.CreateCredentialsResponse
		expectErr error
	}{
		{
			name: "OK",

			request: &services.CreateCredentialsRequest{
				Email:                  "user@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			shouldCallCreateCredentialsDAO: true,
			createCredentialsDAOResponse: &entities.Credential{
				ID:                     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Email:                  "user@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
				CreatedAt:              time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &services.CreateCredentialsResponse{
				ID:                     "00000000-0000-0000-0000-000000000004",
				Email:                  "user@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
				CreatedAt:              time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "OK/Minimal",

			request: &services.CreateCredentialsRequest{
				Email: "user@gmail.com",
				Role:  entities.RoleNone,
			},

			shouldCallCreateCredentialsDAO: true,
			createCredentialsDAOResponse: &entities.Credential{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Email:     "user@gmail.com",
				Role:      entities.RoleNone,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &services.CreateCredentialsResponse{
				ID:        "00000000-0000-0000-0000-000000000004",
				Email:     "user@gmail.com",
				Role:      entities.RoleNone,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "DAO/Error",

			request: &services.CreateCredentialsRequest{
				Email:                  "user@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			shouldCallCreateCredentialsDAO: true,
			createCredentialsDAOError:      errors.New("uwups"),

			expectErr: services.ErrCreateCredentials,
		},
		{
			name: "Invalid/EmailMissing",

			request: &services.CreateCredentialsRequest{
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			expectErr: services.ErrInvalidCreateCredentialsRequest,
		},
		{
			name: "Invalid/EmailTooLong",

			request: &services.CreateCredentialsRequest{
				Email:                  strings.Repeat("a", 256) + "@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			expectErr: services.ErrInvalidCreateCredentialsRequest,
		},
		{
			name: "Invalid/EmailInvalid",

			request: &services.CreateCredentialsRequest{
				Email:                  "fake-email",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			expectErr: services.ErrInvalidCreateCredentialsRequest,
		},
		{
			name: "Invalid/Role",

			request: &services.CreateCredentialsRequest{
				Email:                  "user@gmail.com",
				Role:                   entities.Role("fake-role"),
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
			},

			expectErr: services.ErrInvalidCreateCredentialsRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			createCredentialsDAO := daomocks.NewMockCreateCredentials(t)

			if testCase.shouldCallCreateCredentialsDAO {
				createCredentialsDAO.
					On(
						"Exec",
						context.Background(),
						mock.MatchedBy(func(id uuid.UUID) bool { return id != uuid.Nil }),
						mock.MatchedBy(func(at time.Time) bool { return at.Unix() > 0 }),
						&dao.CreateCredentialsRequest{
							Email:                  testCase.request.Email,
							Role:                   testCase.request.Role,
							EmailValidationTokenID: testCase.request.EmailValidationTokenID,
							PasswordTokenID:        testCase.request.PasswordTokenID,
							ResetPasswordTokenID:   testCase.request.ResetPasswordTokenID,
						},
					).
					Return(testCase.createCredentialsDAOResponse, testCase.createCredentialsDAOError)
			}

			service := services.NewCreateCredentials(createCredentialsDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			createCredentialsDAO.AssertExpectations(t)
		})
	}
}
