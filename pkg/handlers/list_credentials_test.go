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

func TestListCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.ListServiceExecRequest

		serviceResp *services.ListCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.ListServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.ListServiceExecRequest{
				Ids: []string{"id-1", "id-2", "id-3"},
			},

			serviceResp: &services.ListCredentialsResponse{
				Credentials: []*services.ListCredentialsResponseCredential{
					{
						ID:                            "00000000-0000-0000-0000-000000000001",
						Email:                         "email-1",
						Role:                          entities.RoleCore,
						EmailValidationTokenID:        "email-validation-token-id-1",
						PendingEmailValidationTokenID: "pending-email-validation-token-id-1",
						PasswordTokenID:               "password-token-id-1",
						ResetPasswordTokenID:          "reset-password-token-id-1",
						CreatedAt:                     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:                     lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        "00000000-0000-0000-0000-000000000003",
						Email:     "email-3",
						Role:      entities.RoleNone,
						CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: &credentialsv1.ListServiceExecResponse{
				Credentials: []*credentialsv1.ListServiceExecResponseElement{
					{
						Id:                            "00000000-0000-0000-0000-000000000001",
						Email:                         "email-1",
						Role:                          commonv1.UserRole_USER_ROLE_CORE,
						EmailValidationTokenId:        "email-validation-token-id-1",
						PendingEmailValidationTokenId: "pending-email-validation-token-id-1",
						PasswordTokenId:               "password-token-id-1",
						ResetPasswordTokenId:          "reset-password-token-id-1",
						CreatedAt:                     timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
						UpdatedAt:                     timestamppb.New(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						Id:        "00000000-0000-0000-0000-000000000003",
						Email:     "email-3",
						Role:      commonv1.UserRole_USER_ROLE_UNSPECIFIED,
						CreatedAt: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
		},
		{
			name: "InvalidArgument",

			request: &credentialsv1.ListServiceExecRequest{
				Ids: []string{"id-1", "id-2", "id-3"},
			},

			serviceErr: services.ErrInvalidListCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",

			request: &credentialsv1.ListServiceExecRequest{
				Ids: []string{"id-1", "id-2", "id-3"},
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockListCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.ListCredentialsRequest{
					IDs: testCase.request.GetIds(),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.ListCredentialsServiceName, mock.Anything)

			handler := handlers.NewListCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
