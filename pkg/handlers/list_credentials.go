package handlers

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers/adapters"

	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

const ListCredentialsServiceName = "list_credentials"

type ListCredentials interface {
	credentialsv1grpc.ListServiceServer
}

type listCredentialsImpl struct {
	service services.ListCredentials
}

var handleListCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidListCredentialsRequest, codes.InvalidArgument).
	Handle

func credentialToListElementProto(
	item *services.ListCredentialsResponseCredential, _ int,
) *credentialsv1.ListServiceExecResponseElement {
	return &credentialsv1.ListServiceExecResponseElement{
		Id:                            item.ID,
		Email:                         item.Email,
		Role:                          entities.RoleConverter.ToProto(item.Role),
		EmailValidationTokenId:        item.EmailValidationTokenID,
		PendingEmailValidationTokenId: item.PendingEmailValidationTokenID,
		PasswordTokenId:               item.PasswordTokenID,
		ResetPasswordTokenId:          item.ResetPasswordTokenID,
		CreatedAt:                     timestamppb.New(item.CreatedAt),
		UpdatedAt:                     grpc.TimestampOptional(item.UpdatedAt),
	}
}

func (handler *listCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.ListServiceExecRequest,
) (*credentialsv1.ListServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.ListCredentialsRequest{IDs: request.GetIds()})
	if err != nil {
		return nil, handleListCredentialsError(err)
	}

	elements := lo.Map(res.Credentials, credentialToListElementProto)

	return &credentialsv1.ListServiceExecResponse{Credentials: elements}, nil
}

func NewListCredentials(service services.ListCredentials, logger adapters.GRPC) ListCredentials {
	handler := &listCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(ListCredentialsServiceName, handler, logger)
}
