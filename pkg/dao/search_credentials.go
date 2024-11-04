package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"

	"github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-credentials/pkg/entities"
)

var sortedRolesString = strings.Join(
	lo.Reduce(
		entities.RolesSorted,
		func(agg []string, item entities.Role, _ int) []string { return append(agg, "'"+item.String()+"'") },
		[]string{},
	),
	",",
)

type SearchCredentialsRequest struct {
	Limit         int
	Offset        int
	Sort          entities.SortCredentials
	SortDirection database.SortDirection
	Emails        []string
	Roles         []entities.Role
}

type SearchCredentials interface {
	Exec(ctx context.Context, request *SearchCredentialsRequest) (uuid.UUIDs, error)
}

type searchCredentialsImpl struct {
	database bun.IDB
}

func (dao *searchCredentialsImpl) Exec(ctx context.Context, request *SearchCredentialsRequest) (uuid.UUIDs, error) {
	credentials := make([]*entities.Credential, 0)

	query := dao.database.
		NewSelect().
		Model(&credentials).
		Column("id").
		Limit(request.Limit).
		Offset(request.Offset)

	// Only apply sorting direction if a sort value is present. Otherwise, ignore it and use default sorting.
	if request.Sort != entities.SortCredentialsNone {
		direction := lo.Switch[database.SortDirection, string](request.SortDirection).
			Case(database.SortDirectionAsc, "ASC").
			Case(database.SortDirectionDesc, "DESC").
			Default("ASC")

		sort := lo.Switch[entities.SortCredentials, string](request.Sort).
			Case(entities.SortCredentialsEmail, "credentials.email").
			Case(
				entities.SortCredentialsRole,
				fmt.Sprintf("array_position(array[%s],credentials.role::text)", sortedRolesString),
			).
			Case(entities.SortCredentialsCreatedAt, "credentials.created_at").
			Case(entities.SortCredentialsUpdatedAt, "credentials.updated_at").
			Default("credentials.email")

		query = query.OrderExpr(sort + " " + direction)
	} else {
		query = query.Order("credentials.email ASC")
	}

	if len(request.Emails) > 1 {
		query = query.Where("email IN (?)", bun.In(request.Emails))
	} else if len(request.Emails) == 1 {
		query = query.Where("email = ?", request.Emails[0])
	}

	if len(request.Roles) > 1 {
		query = query.Where("role IN (?)", bun.In(request.Roles))
	} else if len(request.Roles) == 1 {
		query = query.Where("role = ?", request.Roles[0])
	}

	err := query.Scan(ctx, &credentials)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	ids := lo.Map(credentials, func(item *entities.Credential, _ int) uuid.UUID {
		return item.ID
	})

	return ids, nil
}

func NewSearchCredentials(database bun.IDB) SearchCredentials {
	return &searchCredentialsImpl{database: database}
}
