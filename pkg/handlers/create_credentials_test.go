package handlers_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestCreateCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.CreateServiceExecRequest

		serviceResp *services.CreateCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.CreateServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.CreateServiceExecRequest{
				Email:                  "email",
				Role:                   commonv1.UserRole_USER_ROLE_EARLY_ACCESS_PROGRAM,
				EmailValidationTokenId: "email-validation",
				PasswordTokenId:        "password",
				ResetPasswordTokenId:   "reset-password",
			},

			serviceResp: &services.CreateCredentialsResponse{
				ID:                     "00000000-0000-0000-0000-000000000004",
				Email:                  "user@gmail.com",
				Role:                   entities.RoleAdmin,
				EmailValidationTokenID: "00000000-0000-0000-0000-000000000001",
				PasswordTokenID:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenID:   "00000000-0000-0000-0000-000000000003",
				CreatedAt:              time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &credentialsv1.CreateServiceExecResponse{
				Id:                     "00000000-0000-0000-0000-000000000004",
				Email:                  "user@gmail.com",
				Role:                   commonv1.UserRole_USER_ROLE_ADMIN,
				EmailValidationTokenId: "00000000-0000-0000-0000-000000000001",
				PasswordTokenId:        "00000000-0000-0000-0000-000000000002",
				ResetPasswordTokenId:   "00000000-0000-0000-0000-000000000003",
				CreatedAt:              timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "InvalidRequest",

			request: &credentialsv1.CreateServiceExecRequest{
				Email:                  "email",
				Role:                   commonv1.UserRole_USER_ROLE_EARLY_ACCESS_PROGRAM,
				EmailValidationTokenId: "email-validation",
				PasswordTokenId:        "password",
				ResetPasswordTokenId:   "reset-password",
			},

			serviceErr: services.ErrInvalidCreateCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "AlreadyExists",

			request: &credentialsv1.CreateServiceExecRequest{
				Email:                  "email",
				Role:                   commonv1.UserRole_USER_ROLE_EARLY_ACCESS_PROGRAM,
				EmailValidationTokenId: "email-validation",
				PasswordTokenId:        "password",
				ResetPasswordTokenId:   "reset-password",
			},

			serviceErr: dao.ErrCredentialsAlreadyExist,

			expectCode: codes.AlreadyExists,
		},
		{
			name: "InternalError",

			request: &credentialsv1.CreateServiceExecRequest{
				Email:                  "email",
				Role:                   commonv1.UserRole_USER_ROLE_EARLY_ACCESS_PROGRAM,
				EmailValidationTokenId: "email-validation",
				PasswordTokenId:        "password",
				ResetPasswordTokenId:   "reset-password",
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockCreateCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.CreateCredentialsRequest{
					Email:                  testCase.request.GetEmail(),
					Role:                   entities.RoleConverter.FromProto(testCase.request.GetRole()),
					EmailValidationTokenID: testCase.request.GetEmailValidationTokenId(),
					PasswordTokenID:        testCase.request.GetPasswordTokenId(),
					ResetPasswordTokenID:   testCase.request.GetResetPasswordTokenId(),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.CreateCredentialsServiceName, mock.Anything)

			handler := handlers.NewCreateCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
