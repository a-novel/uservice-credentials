package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/testutils"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/handlers"
	"github.com/a-novel/uservice-credentials/pkg/services"
	servicesmocks "github.com/a-novel/uservice-credentials/pkg/services/mocks"
)

func TestExistsCredentials(t *testing.T) {
	testCases := []struct {
		name string

		request *credentialsv1.ExistsServiceExecRequest

		serviceResp *services.ExistsCredentialsResponse
		serviceErr  error

		expect     *credentialsv1.ExistsServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			request: &credentialsv1.ExistsServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceResp: &services.ExistsCredentialsResponse{
				Exists: true,
			},

			expect: &credentialsv1.ExistsServiceExecResponse{
				Exists: true,
			},
		},
		{
			name: "InvalidRequest",

			request: &credentialsv1.ExistsServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: services.ErrInvalidExistsCredentialsRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "NotFound",

			request: &credentialsv1.ExistsServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: dao.ErrCredentialsNotFound,

			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",

			request: &credentialsv1.ExistsServiceExecRequest{
				Id:    "00000000-0000-0000-0000-000000000004",
				Email: "email",
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockExistsCredentials(t)
			logger := adaptersmocks.NewMockGRPC(t)

			service.
				On("Exec", context.Background(), &services.ExistsCredentialsRequest{
					ID:    testCase.request.GetId(),
					Email: testCase.request.GetEmail(),
				}).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.ExistsCredentialsServiceName, mock.Anything)

			handler := handlers.NewExistsCredentials(service, logger)
			resp, err := handler.Exec(context.Background(), testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
