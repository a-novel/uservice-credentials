package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"
	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	credentialsv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/credentials/v1"

	anovelgrpc "github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/testutils"
)

func init() {
	go main()
}

var servicesToTest = []string{
	"create",
	"get",
	"list",
	"search",
	"update",
}

func TestIntegrationHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	// Create the RPC client.
	pool := anovelgrpc.NewConnPool()
	conn, err := pool.Open("0.0.0.0", 8080, anovelgrpc.ProtocolHTTP)
	require.NoError(t, err)

	healthClient := healthpb.NewHealthClient(conn)

	testutils.WaitConn(t, conn)

	for _, service := range servicesToTest {
		res, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{Service: service})
		require.NoError(t, err)
		require.Equal(t, healthpb.HealthCheckResponse_SERVING, res.Status)
	}
}

func TestIntegrationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	// Create the RPC client.
	pool := anovelgrpc.NewConnPool()
	conn, err := pool.Open("0.0.0.0", 8080, anovelgrpc.ProtocolHTTP)
	require.NoError(t, err)

	testutils.WaitConn(t, conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createCredentialsClient := credentialsv1grpc.NewCreateServiceClient(conn)
	getCredentialsClient := credentialsv1grpc.NewGetServiceClient(conn)
	updateCredentialsClient := credentialsv1grpc.NewUpdateServiceClient(conn)

	// Create credentials
	createResp, err := createCredentialsClient.Exec(ctx, &credentialsv1.CreateServiceExecRequest{
		Email:                  "user@gmail.com",
		Role:                   commonv1.UserRole_USER_ROLE_CORE,
		EmailValidationTokenId: "email-token",
		PasswordTokenId:        "password-token",
		ResetPasswordTokenId:   "reset-password-token",
	})
	require.NoError(t, err)
	require.NotEmpty(t, createResp.Id)
	require.Equal(t, "user@gmail.com", createResp.Email)
	require.Equal(t, commonv1.UserRole_USER_ROLE_CORE, createResp.Role)
	require.Equal(t, "email-token", createResp.EmailValidationTokenId)
	require.Equal(t, "password-token", createResp.PasswordTokenId)
	require.Equal(t, "reset-password-token", createResp.ResetPasswordTokenId)

	// Get credentials (by email)
	getResp, err := getCredentialsClient.Exec(ctx, &credentialsv1.GetServiceExecRequest{
		Email: "user@gmail.com",
	})
	require.NoError(t, err)
	require.Equal(t, createResp.Id, getResp.Id)
	require.Equal(t, createResp.Email, getResp.Email)
	require.Equal(t, createResp.Role, getResp.Role)
	require.Equal(t, createResp.EmailValidationTokenId, getResp.EmailValidationTokenId)
	require.Equal(t, createResp.PasswordTokenId, getResp.PasswordTokenId)
	require.Equal(t, createResp.ResetPasswordTokenId, getResp.ResetPasswordTokenId)
	require.Equal(t, createResp.CreatedAt, getResp.CreatedAt)

	// Get credentials (by id)
	getResp, err = getCredentialsClient.Exec(ctx, &credentialsv1.GetServiceExecRequest{
		Id: createResp.Id,
	})
	require.NoError(t, err)
	require.Equal(t, createResp.Id, getResp.Id)
	require.Equal(t, createResp.Email, getResp.Email)
	require.Equal(t, createResp.Role, getResp.Role)
	require.Equal(t, createResp.EmailValidationTokenId, getResp.EmailValidationTokenId)
	require.Equal(t, createResp.PasswordTokenId, getResp.PasswordTokenId)
	require.Equal(t, createResp.ResetPasswordTokenId, getResp.ResetPasswordTokenId)
	require.Equal(t, createResp.CreatedAt, getResp.CreatedAt)

	// Update credentials
	updateResp, err := updateCredentialsClient.Exec(ctx, &credentialsv1.UpdateServiceExecRequest{
		Id:                            createResp.Id,
		Email:                         "user-alt@gmail.com",
		Role:                          commonv1.UserRole_USER_ROLE_ADMIN,
		EmailValidationTokenId:        "",
		PendingEmailValidationTokenId: "pending-email-token",
		PasswordTokenId:               "",
		ResetPasswordTokenId:          "",
	})
	require.NoError(t, err)
	require.Equal(t, createResp.Id, updateResp.Id)
	require.Equal(t, "user-alt@gmail.com", updateResp.Email)
	require.Equal(t, commonv1.UserRole_USER_ROLE_ADMIN, updateResp.Role)
	require.Empty(t, updateResp.EmailValidationTokenId)
	require.Equal(t, "pending-email-token", updateResp.PendingEmailValidationTokenId)
	require.Empty(t, updateResp.PasswordTokenId)
	require.Empty(t, updateResp.ResetPasswordTokenId)
	require.Equal(t, createResp.CreatedAt, updateResp.CreatedAt)
}
