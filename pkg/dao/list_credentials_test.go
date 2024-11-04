package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	anoveldb "github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-credentials/migrations"
	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
)

func TestListCredentials(t *testing.T) {
	fixtures := []interface{}{
		&entities.Credential{
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
		&entities.Credential{
			ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Email:                         "email-2",
			Role:                          entities.RoleAdmin,
			EmailValidationTokenID:        "email-validation-token-id-2",
			PendingEmailValidationTokenID: "pending-email-validation-token-id-2",
			PasswordTokenID:               "password-token-id-2",
			ResetPasswordTokenID:          "reset-password-token-id-2",
			CreatedAt:                     time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:                     lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
		},
		&entities.Credential{
			ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Email:                         "email-3",
			Role:                          entities.RoleNone,
			EmailValidationTokenID:        "email-validation-token-id-3",
			PendingEmailValidationTokenID: "pending-email-validation-token-id-3",
			PasswordTokenID:               "password-token-id-3",
			ResetPasswordTokenID:          "reset-password-token-id-3",
			CreatedAt:                     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:                     lo.ToPtr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	testCases := []struct {
		name string

		ids []uuid.UUID

		expect    []*entities.Credential
		expectErr error
	}{
		{
			name: "List",

			ids: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},

			expect: []*entities.Credential{
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
					ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Email:                         "email-3",
					Role:                          entities.RoleNone,
					EmailValidationTokenID:        "email-validation-token-id-3",
					PendingEmailValidationTokenID: "pending-email-validation-token-id-3",
					PasswordTokenID:               "password-token-id-3",
					ResetPasswordTokenID:          "reset-password-token-id-3",
					CreatedAt:                     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:                     lo.ToPtr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "IgnoreMissingIDs",

			ids: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			},

			expect: []*entities.Credential{
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
					ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Email:                         "email-3",
					Role:                          entities.RoleNone,
					EmailValidationTokenID:        "email-validation-token-id-3",
					PendingEmailValidationTokenID: "pending-email-validation-token-id-3",
					PasswordTokenID:               "password-token-id-3",
					ResetPasswordTokenID:          "reset-password-token-id-3",
					CreatedAt:                     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:                     lo.ToPtr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "NoResults",

			ids: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			},

			expect: []*entities.Credential{},
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	transaction := anoveldb.BeginTestTX(database, fixtures)
	defer anoveldb.RollbackTestTX(transaction)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			listCredentialsDAO := dao.NewListCredentials(transaction)

			credential, err := listCredentialsDAO.Exec(context.Background(), testCase.ids)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, credential)
		})
	}
}
