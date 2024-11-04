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

func TestUpdateCredentials(t *testing.T) {
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
			UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
	}

	testCases := []struct {
		name string

		id   uuid.UUID
		now  time.Time
		data *dao.UpdateCredentialsRequest

		expect    *entities.Credential
		expectErr error
	}{
		{
			name: "Update",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			data: &dao.UpdateCredentialsRequest{
				Email:                         "email-2",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "new-email-validation-token-id",
				PendingEmailValidationTokenID: "new-pending-email-validation-token-id",
				PasswordTokenID:               "new-password-token-id",
				ResetPasswordTokenID:          "new-reset-password-token-id",
			},

			expect: &entities.Credential{
				ID:                            uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Email:                         "email-2",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "new-email-validation-token-id",
				PendingEmailValidationTokenID: "new-pending-email-validation-token-id",
				PasswordTokenID:               "new-password-token-id",
				ResetPasswordTokenID:          "new-reset-password-token-id",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "Update/Clear",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			data: &dao.UpdateCredentialsRequest{
				Email: "email-2",
				Role:  entities.RoleAdmin,
			},

			expect: &entities.Credential{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Email:     "email-2",
				Role:      entities.RoleAdmin,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "NotFound",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			data: &dao.UpdateCredentialsRequest{
				Email:                         "email-2",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "new-email-validation-token-id",
				PendingEmailValidationTokenID: "new-pending-email-validation-token-id",
				PasswordTokenID:               "new-password-token-id",
				ResetPasswordTokenID:          "new-reset-password-token-id",
			},

			expectErr: dao.ErrCredentialsNotFound,
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			transaction := anoveldb.BeginTestTX(database, fixtures)
			defer anoveldb.RollbackTestTX(transaction)

			updateCredentialsDAO := dao.NewUpdateCredentials(transaction)

			beat, err := updateCredentialsDAO.Exec(context.Background(), testCase.id, testCase.now, testCase.data)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, beat)
		})
	}
}
