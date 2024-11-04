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

const UpdateCredentialsServiceName = "update_credentials" //nolint:gosec

type UpdateCredentials interface {
	credentialsv1grpc.UpdateServiceServer
}

type updateCredentialsImpl struct {
	service services.UpdateCredentials
}

var handleUpdateCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidUpdateCredentialsRequest, codes.InvalidArgument).
	Is(dao.ErrCredentialsNotFound, codes.NotFound).
	Handle

func (handler *updateCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.UpdateServiceExecRequest,
) (*credentialsv1.UpdateServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.UpdateCredentialsRequest{
		ID:                            request.GetId(),
		Email:                         request.GetEmail(),
		Role:                          entities.RoleConverter.FromProto(request.GetRole()),
		EmailValidationTokenID:        request.GetEmailValidationTokenId(),
		PendingEmailValidationTokenID: request.GetPendingEmailValidationTokenId(),
		PasswordTokenID:               request.GetPasswordTokenId(),
		ResetPasswordTokenID:          request.GetResetPasswordTokenId(),
	})
	if err != nil {
		return nil, handleUpdateCredentialsError(err)
	}

	return &credentialsv1.UpdateServiceExecResponse{
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

func NewUpdateCredentials(service services.UpdateCredentials, logger adapters.GRPC) UpdateCredentials {
	handler := &updateCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(UpdateCredentialsServiceName, handler, logger)
}
