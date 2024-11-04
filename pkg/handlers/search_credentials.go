package handlers

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"
	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	"github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers/adapters"

	"github.com/a-novel/uservice-credentials/pkg/entities"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

const SearchCredentialsServiceName = "search_credentials"

type SearchCredentials interface {
	credentialsv1grpc.SearchServiceServer
}

type searchCredentialsImpl struct {
	service services.SearchCredentials
}

var handleSearchCredentialsError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidSearchCredentialsRequest, codes.InvalidArgument).
	Handle

func (handler *searchCredentialsImpl) Exec(
	ctx context.Context, request *credentialsv1.SearchServiceExecRequest,
) (*credentialsv1.SearchServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.SearchCredentialsRequest{
		Limit:         int(request.GetPagination().GetLimit()),
		Offset:        int(request.GetPagination().GetOffset()),
		Sort:          entities.SortCredentialsConverter.FromProto(request.GetOrderBy()),
		SortDirection: grpc.SortDirectionConverter.FromProto(request.GetOrderDirection()),
		Emails:        request.GetEmails(),
		Roles: lo.Map(request.GetRoles(), func(item commonv1.UserRole, _ int) entities.Role {
			return entities.RoleConverter.FromProto(item)
		}),
	})
	if err != nil {
		return nil, handleSearchCredentialsError(err)
	}

	return &credentialsv1.SearchServiceExecResponse{Ids: res.IDs}, nil
}

func NewSearchCredentials(service services.SearchCredentials, logger adapters.GRPC) SearchCredentials {
	handler := &searchCredentialsImpl{service: service}
	return grpc.ServiceWithMetrics(SearchCredentialsServiceName, handler, logger)
}
