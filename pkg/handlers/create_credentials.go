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

const CreateCredentialsServiceName = "create_credentials" //nolint:gosec

type CreateCredentials interface {
	credentialsv1grpc.CreateServiceServer
}

type createCredentialsImpl struct {
	service services.CreateCredentials
}

var handleCreateCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidCreateCredentialsRequest, codes.InvalidArgument).
	Is(dao.ErrCredentialsAlreadyExist, codes.AlreadyExists).
	Handle

func (handler *createCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.CreateServiceExecRequest,
) (*credentialsv1.CreateServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.CreateCredentialsRequest{
		Email:                  request.GetEmail(),
		Role:                   entities.RoleConverter.FromProto(request.GetRole()),
		EmailValidationTokenID: request.GetEmailValidationTokenId(),
		PasswordTokenID:        request.GetPasswordTokenId(),
		ResetPasswordTokenID:   request.GetResetPasswordTokenId(),
	})
	if err != nil {
		return nil, handleCreateCredentialsError(err)
	}

	return &credentialsv1.CreateServiceExecResponse{
		Id:                     res.ID,
		Email:                  res.Email,
		Role:                   entities.RoleConverter.ToProto(res.Role),
		EmailValidationTokenId: res.EmailValidationTokenID,
		PasswordTokenId:        res.PasswordTokenID,
		ResetPasswordTokenId:   res.ResetPasswordTokenID,
		CreatedAt:              timestamppb.New(res.CreatedAt),
	}, nil
}

func NewCreateCredentials(service services.CreateCredentials, logger adapters.GRPC) CreateCredentials {
	handler := &createCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(CreateCredentialsServiceName, handler, logger)
}
