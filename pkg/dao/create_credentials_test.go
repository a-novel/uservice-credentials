package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	anoveldb "github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-credentials/migrations"
	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

func TestCreateCredentials(t *testing.T) {
	fixtures := []interface{}{
		&entities.Credential{
			ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:                         "email-1",
			Role:                          entities.RoleCore,
			EmailValidationTokenID:        "email-validation-token-id",
			PendingEmailValidationTokenID: "pending-email-validation-token-id",
			PasswordTokenID:               "password-token-id",
			ResetPasswordTokenID:          "reset-password-token-id",
			CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name string

		id  uuid.UUID
		now time.Time

		request *dao.CreateCredentialsRequest

		expect    *entities.Credential
		expectErr error
	}{
		{
			name: "Create",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.CreateCredentialsRequest{
				Email:                  "email-2",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "email-validation-token-id",
				PasswordTokenID:        "password-token-id",
				ResetPasswordTokenID:   "reset-password-token-id",
			},

			expect: &entities.Credential{
				ID:                     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Email:                  "email-2",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "email-validation-token-id",
				PasswordTokenID:        "password-token-id",
				ResetPasswordTokenID:   "reset-password-token-id",
				CreatedAt:              time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Create/Minimal",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.CreateCredentialsRequest{
				Email: "email-2",
			},

			expect: &entities.Credential{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Email:     "email-2",
				Role:      entities.RoleNone,
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Create/EmailAlreadyExists",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.CreateCredentialsRequest{
				Email: "email-1",
			},

			expectErr: dao.ErrCredentialsAlreadyExist,
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			transaction := anoveldb.BeginTestTX(database, fixtures)
			defer anoveldb.RollbackTestTX(transaction)

			createCredentialsDAO := dao.NewCreateCredentials(transaction)

			credential, err := createCredentialsDAO.
				Exec(context.Background(), testCase.id, testCase.now, testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, credential)
		})
	}
}
