package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers/adapters"

	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

const GetCredentialsServiceName = "get_credentials"

type GetCredentials interface {
	credentialsv1grpc.GetServiceServer
}

type getCredentialsImpl struct {
	service services.GetCredentials
}

var handleGetCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidGetCredentialsRequest, codes.InvalidArgument).
	Is(dao.ErrCredentialsNotFound, codes.NotFound).
	Handle

func (handler *getCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.GetServiceExecRequest,
) (*credentialsv1.GetServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.GetCredentialsRequest{
		ID:    request.GetId(),
		Email: request.GetEmail(),
	})
	if err != nil {
		return nil, handleGetCredentialsError(err)
	}

	return &credentialsv1.GetServiceExecResponse{
		Id:                            res.ID,
		Email:                         res.Email,
		Role:                          entities.RoleConverter.ToProto(res.Role),
		EmailValidationTokenId:        res.EmailValidationTokenID,
		PendingEmailValidationTokenId: res.PendingEmailValidationTokenID,
		PasswordTokenId:               res.PasswordTokenID,
		ResetPasswordTokenId:          res.ResetPasswordTokenID,
		CreatedAt:                     timestamppb.New(res.CreatedAt),
		UpdatedAt:                     grpc.TimestampOptional(res.UpdatedAt),
	}, nil
}

func NewGetCredentials(service services.GetCredentials, logger adapters.GRPC) GetCredentials {
	handler := &getCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(GetCredentialsServiceName, handler, logger)
}
