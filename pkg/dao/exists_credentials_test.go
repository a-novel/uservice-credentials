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

func TestExistsCredentials(t *testing.T) {
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
			UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	testCases := []struct {
		name string

		request *dao.ExistsCredentialsRequest

		expect    bool
		expectErr error
	}{
		{
			name: "Exists/Email",

			request: &dao.ExistsCredentialsRequest{
				Email: "email-1",
			},

			expect: true,
		},
		{
			name: "Exists/ID",

			request: &dao.ExistsCredentialsRequest{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},

			expect: true,
		},
		{
			name: "Exists/NotFound",

			request: &dao.ExistsCredentialsRequest{
				Email: "email-2",
			},

			expect: false,
		},
		{
			name: "Exists/NoParameters",

			request: &dao.ExistsCredentialsRequest{},

			expectErr: dao.ErrCredentialsNotFound,
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	transaction := anoveldb.BeginTestTX(database, fixtures)
	defer anoveldb.RollbackTestTX(transaction)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			existsCredentialsDAO := dao.NewExistsCredentials(transaction)

			credential, err := existsCredentialsDAO.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, credential)
		})
	}
}
