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

func TestCredentialsSearch(t *testing.T) {
	fixtures := []interface{}{
		// Order by email: Credentials 1, Credentials 2, Credentials 3
		// Order by role: Credentials 2, Credentials 3, Credentials 1
		// Order by created_at: Credentials 3, Credentials 2, Credentials 1
		// Order by updated_at: Credentials 3, Credentials 1, Credentials 2
		// Insertion order: Credentials 2, Credentials 1, Credentials 3

		&entities.Credential{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Email:     "email_2",
			Role:      entities.RoleEarlyAccessProgram,
			CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: lo.ToPtr(time.Date(2021, 4, 2, 0, 0, 0, 0, time.UTC)),
		},
		&entities.Credential{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:     "email_1",
			Role:      entities.RoleCore,
			CreatedAt: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: lo.ToPtr(time.Date(2021, 3, 2, 0, 0, 0, 0, time.UTC)),
		},
		&entities.Credential{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Email:     "email_3",
			Role:      entities.RoleAdmin,
			CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
	}

	testCases := []struct {
		name string

		request *dao.SearchCredentialsRequest

		expect    uuid.UUIDs
		expectErr error
	}{
		// Base.
		{
			name: "Base",

			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
			},

			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},

		// Limit.
		{
			name: "LimitTooLow",
			request: &dao.SearchCredentialsRequest{
				Limit:  2,
				Offset: 0,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		},
		{
			name: "LimitTooHigh",
			request: &dao.SearchCredentialsRequest{
				Limit:  10,
				Offset: 0,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},

		// Offset.
		{
			name: "Offset",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 1,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "OffsetTooHigh",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 10,
			},
			expect: uuid.UUIDs{},
		},

		// Sort: email
		{
			name: "Email",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Sort:   entities.SortCredentialsEmail,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "EmailAsc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionAsc,
				Sort:          entities.SortCredentialsEmail,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "EmailDesc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionDesc,
				Sort:          entities.SortCredentialsEmail,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},

		// Sort: role
		{
			name: "Role",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Sort:   entities.SortCredentialsRole,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},
		{
			name: "RoleAsc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionAsc,
				Sort:          entities.SortCredentialsRole,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},
		{
			name: "RoleDesc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionDesc,
				Sort:          entities.SortCredentialsRole,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		},

		// Sort: created_at
		{
			name: "CreatedAt",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Sort:   entities.SortCredentialsCreatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},
		{
			name: "CreatedAtAsc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionAsc,
				Sort:          entities.SortCredentialsCreatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},
		{
			name: "CreatedAtDesc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionDesc,
				Sort:          entities.SortCredentialsCreatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},

		// Sort: updated_at
		{
			name: "UpdatedAt",
			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Sort:   entities.SortCredentialsUpdatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		},
		{
			name: "UpdatedAtAsc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionAsc,
				Sort:          entities.SortCredentialsUpdatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		},
		{
			name: "UpdatedAtDesc",
			request: &dao.SearchCredentialsRequest{
				Limit:         3,
				Offset:        0,
				SortDirection: anoveldb.SortDirectionDesc,
				Sort:          entities.SortCredentialsUpdatedAt,
			},
			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},

		// Filter: email
		{
			name: "Filter/Email",

			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Emails: []string{"email_1", "email_2"},
			},

			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			},
		},
		{
			name: "Filter/Email/Simple",

			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Emails: []string{"email_1"},
			},

			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},

		// Filter: roles
		{
			name: "Filter/Roles",

			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Roles:  []entities.Role{entities.RoleCore, entities.RoleAdmin},
			},

			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "Filter/Roles/simple",

			request: &dao.SearchCredentialsRequest{
				Limit:  3,
				Offset: 0,
				Roles:  []entities.Role{entities.RoleAdmin},
			},

			expect: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	transaction := anoveldb.BeginTestTX(database, fixtures)
	defer anoveldb.RollbackTestTX(transaction)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			searchCredentialsDAO := dao.NewSearchCredentials(transaction)

			credential, err := searchCredentialsDAO.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, credential)
		})
	}
}
