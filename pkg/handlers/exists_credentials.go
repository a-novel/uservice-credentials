package handlers

import (
	"context"

	"google.golang.org/grpc/codes"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers/adapters"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

const ExistsCredentialsServiceName = "exists_credentials" //nolint:gosec

type ExistsCredentials interface {
	credentialsv1grpc.ExistsServiceServer
}

type existsCredentialsImpl struct {
	service services.ExistsCredentials
}

var handleExistsCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidExistsCredentialsRequest, codes.InvalidArgument).
	Is(dao.ErrCredentialsNotFound, codes.NotFound).
	Handle

func (handler *existsCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.ExistsServiceExecRequest,
) (*credentialsv1.ExistsServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.ExistsCredentialsRequest{
		ID:    request.GetId(),
		Email: request.GetEmail(),
	})
	if err != nil {
		return nil, handleExistsCredentialsError(err)
	}

	return &credentialsv1.ExistsServiceExecResponse{
		Exists: res.Exists,
	}, nil
}

func NewExistsCredentials(service services.ExistsCredentials, logger adapters.GRPC) ExistsCredentials {
	handler := &existsCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(ExistsCredentialsServiceName, handler, logger)
}
