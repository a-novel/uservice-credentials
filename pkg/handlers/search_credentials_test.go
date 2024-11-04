package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/grpc"
	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/testutils"

	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/handlers"
	"github.com/a-novel/uservice-credentials/pkg/services"
	servicesmocks "github.com/a-novel/uservice-credentials/pkg/services/mocks"
)

func TestSearchCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.SearchServiceExecRequest

		serviceResp *services.SearchCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.SearchServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.SearchServiceExecRequest{
				Pagination: &commonv1.Pagination{
					Limit:  10,
					Offset: 2,
				},
				OrderBy:        credentialsv1.Sort_SORT_BY_EMAIL,
				OrderDirection: commonv1.SortDirection_SORT_DIRECTION_ASC,
				Emails:         []string{"email-1", "email-2"},
				Roles:          []commonv1.UserRole{commonv1.UserRole_USER_ROLE_CORE},
			},

			serviceResp: &services.SearchCredentialsResponse{
				IDs: []string{"id-1", "id-2", "id-3"},
			},

			expect: &credentialsv1.SearchServiceExecResponse{
				Ids: []string{"id-1", "id-2", "id-3"},
			},
		},
		{
			name: "InvalidArgument",

			request: &credentialsv1.SearchServiceExecRequest{
				Pagination: &commonv1.Pagination{
					Limit:  10,
					Offset: 2,
				},
				OrderBy:        credentialsv1.Sort_SORT_BY_EMAIL,
				OrderDirection: commonv1.SortDirection_SORT_DIRECTION_ASC,
				Emails:         []string{"email-1", "email-2"},
				Roles:          []commonv1.UserRole{commonv1.UserRole_USER_ROLE_CORE},
			},

			serviceErr: services.ErrInvalidSearchCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",

			request: &credentialsv1.SearchServiceExecRequest{
				Pagination: &commonv1.Pagination{
					Limit:  10,
					Offset: 2,
				},
				OrderBy:        credentialsv1.Sort_SORT_BY_EMAIL,
				OrderDirection: commonv1.SortDirection_SORT_DIRECTION_ASC,
				Emails:         []string{"email-1", "email-2"},
				Roles:          []commonv1.UserRole{commonv1.UserRole_USER_ROLE_CORE},
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockSearchCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.SearchCredentialsRequest{
					Limit:         int(testCase.request.GetPagination().GetLimit()),
					Offset:        int(testCase.request.GetPagination().GetOffset()),
					Sort:          entities.SortCredentialsConverter.FromProto(testCase.request.GetOrderBy()),
					SortDirection: grpc.SortDirectionConverter.FromProto(testCase.request.GetOrderDirection()),
					Emails:        testCase.request.GetEmails(),
					Roles: lo.Map(testCase.request.GetRoles(), func(item commonv1.UserRole, _ int) entities.Role {
						return entities.RoleConverter.FromProto(item)
					}),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.SearchCredentialsServiceName, mock.Anything)

			handler := handlers.NewSearchCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
