package handlers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/testutils"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/handlers"
	"github.com/a-novel/uservice-credentials/pkg/services"
	servicesmocks "github.com/a-novel/uservice-credentials/pkg/services/mocks"
)

func TestGetCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.GetServiceExecRequest

		serviceResp *services.GetCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.GetServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.GetServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceResp: &services.GetCredentialsResponse{
				ID:                            "00000000-0000-0000-0000-000000000004",
				Email:                         "email",
				Role:                          entities.RoleAdmin,
				EmailValidationTokenID:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenID: "00000000-0000-0000-0000-000000000005",
				PasswordTokenID:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &credentialsv1.GetServiceExecResponse{
				Id:                            "00000000-0000-0000-0000-000000000004",
				Email:                         "email",
				Role:                          commonv1.UserRole_USER_ROLE_ADMIN,
				EmailValidationTokenId:        "00000000-0000-0000-0000-000000000001",
				PendingEmailValidationTokenId: "00000000-0000-0000-0000-000000000005",
				PasswordTokenId:               "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenId:          "00000000-0000-0000-0000-000000000003",
				CreatedAt:                     timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:                     timestamppb.New(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "InvalidRequest",

			request: &credentialsv1.GetServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: services.ErrInvalidGetCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "NotFound",

			request: &credentialsv1.GetServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: dao.ErrCredentialsNotFound,

			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",

			request: &credentialsv1.GetServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.GetCredentialsRequest{
					ID:    testCase.request.GetId(),
					Email: testCase.request.GetEmail(),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.GetCredentialsServiceName, mock.Anything)

			handler := handlers.NewGetCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
