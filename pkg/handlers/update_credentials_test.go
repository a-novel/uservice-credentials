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

	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/handlers"
	"github.com/a-novel/uservice-credentials/pkg/services"
	servicesmocks "github.com/a-novel/uservice-credentials/pkg/services/mocks"
)

func TestUpdateCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.UpdateServiceExecRequest

		serviceResp *services.UpdateCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.UpdateServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.UpdateServiceExecRequest{
				Id:                            "id",
				Email:                         "email",
				Role:                          commonv1.UserRole_USER_ROLE_CORE,
				EmailValidationTokenId:        "email-validation",
				PendingEmailValidationTokenId: "pending-email-validation",
				PasswordTokenId:               "password",
				ResetPasswordTokenId:          "reset-password",
			},

			serviceResp: &services.UpdateCredentialsResponse{
				ID:                            "id",
				Email:                         "email",
				Role:                          entities.RoleCore,
				EmailValidationTokenID:        "email-validation",
				PendingEmailValidationTokenID: "pending-email-validation",
				PasswordTokenID:               "password",
				ResetPasswordTokenID:          "reset-password",
				CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &credentialsv1.UpdateServiceExecResponse{
				Id:                            "id",
				Email:                         "email",
				Role:                          commonv1.UserRole_USER_ROLE_CORE,
				EmailValidationTokenId:        "email-validation",
				PendingEmailValidationTokenId: "pending-email-validation",
				PasswordTokenId:               "password",
				ResetPasswordTokenId:          "reset-password",
				CreatedAt:                     timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:                     timestamppb.New(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "InvalidArgument",

			request: &credentialsv1.UpdateServiceExecRequest{
				Id:                            "id",
				Email:                         "email",
				Role:                          commonv1.UserRole_USER_ROLE_CORE,
				EmailValidationTokenId:        "email-validation",
				PendingEmailValidationTokenId: "pending-email-validation",
				PasswordTokenId:               "password",
				ResetPasswordTokenId:          "reset-password",
			},

			serviceErr: services.ErrInvalidUpdateCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",

			request: &credentialsv1.UpdateServiceExecRequest{
				Id:                            "id",
				Email:                         "email",
				Role:                          commonv1.UserRole_USER_ROLE_CORE,
				EmailValidationTokenId:        "email-validation",
				PendingEmailValidationTokenId: "pending-email-validation",
				PasswordTokenId:               "password",
				ResetPasswordTokenId:          "reset-password",
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpdateCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.UpdateCredentialsRequest{
					ID:                            testCase.request.GetId(),
					Email:                         testCase.request.GetEmail(),
					Role:                          entities.RoleConverter.FromProto(testCase.request.GetRole()),
					EmailValidationTokenID:        testCase.request.GetEmailValidationTokenId(),
					PendingEmailValidationTokenID: testCase.request.GetPendingEmailValidationTokenId(),
					PasswordTokenID:               testCase.request.GetPasswordTokenId(),
					ResetPasswordTokenID:          testCase.request.GetResetPasswordTokenId(),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.UpdateCredentialsServiceName, mock.Anything)

			handler := handlers.NewUpdateCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
